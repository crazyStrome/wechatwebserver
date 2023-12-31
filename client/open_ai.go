package client

import (
	"context"
	"fmt"
	"time"
	"wechatwebserver/config"

	"github.com/patrickmn/go-cache"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

var openAICli *openai.Client

var openAICache *cache.Cache

func Init() error {
	cfg := openai.DefaultConfig(config.GetConf().OpenAIConfig.Token)
	cfg.BaseURL = config.GetConf().OpenAIConfig.URL
	openAICli = openai.NewClientWithConfig(cfg)

	openAICache = cache.New(
        time.Duration(config.GetConf().OpenAICache.ExpireTime)*time.Second,
		time.Duration(config.GetConf().OpenAICache.CleanTime)*time.Second,
	)
	return nil
}

func Test(ctx context.Context) {
	resp, err := openAICli.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		logrus.Errorf("ChatCompletion error: %v\n", err)
		return
	}

	logrus.Infof(resp.Choices[0].Message.Content)
}

func Talk(ctx context.Context, msgID uint64, reqMsg string) (string, error) {
    cacheMsg, ok := getTalkCache(msgID)
    if ok {
        logrus.Infof("Talk get cache:%v is:%v", msgID, cacheMsg)
        return cacheMsg, nil
    }
    canceled := false
	ch := make(chan string)
    defer close(ch)
	now := time.Now()
	go func() {
		rsp, err := openAICli.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: reqMsg,
				},
			},
		})
		if err != nil {
			logrus.Errorf("Talk:%v err:%v, cost:%v", reqMsg, err, time.Since(now))
			return
		}
        if canceled {
            // 如果已经cancel了，就不往 ch 写了，肯定会 panic
			logrus.Infof("Talk:%v cost:%v, rsp:%v", reqMsg, time.Since(now), rsp.Choices[0].Message.Content)
            // 写到缓存里，下次重试的话可以用
            addTalkCache(msgID, rsp.Choices[0].Message.Content)
            return
		}
		ch <- rsp.Choices[0].Message.Content
	}()
	select {
	case <-ctx.Done():
        logrus.Infof("Talk:%v cost:%v, ctx.done:%v", reqMsg, time.Since(now), ctx.Err())
        canceled = true
		return "", ctx.Err()
	case data := <-ch:
		return data, nil
	}
}

func getTalkCache(msgID uint64) (string, bool) {
    msg := fmt.Sprintf("%v", msgID)
    value, ok := openAICache.Get(msg)
    if !ok {
        return "", false
    }
    ret, ok := value.(string)
    if !ok {
        return "", false
    }
    // 这里能取到的话，后面就不会重试了，直接删掉，节省内存
    openAICache.Delete(msg)
    return ret, true
}

// addTalkCache msg 过长的话，也不会写进去，阈值是 256 字节
func addTalkCache(msgID uint64, msg string) {
    cfg := config.GetConf()
    if len(msg) > int(cfg.OpenAICache.CacheLen) {
        return
    }
    id := fmt.Sprintf("%v", msgID)
    openAICache.Add(id, msg, time.Duration(cfg.OpenAICache.ExpireTime) * time.Second)
}
