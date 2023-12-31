package client

import (
	"context"
	"fmt"
	"sync"
	"time"
	"wechatwebserver/config"

	"github.com/patrickmn/go-cache"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

var openAICli *openai.Client

var (
	openAICache   *cache.Cache
	openAIWorkers map[uint64]chan struct{}
	openAILock    sync.Mutex
)

func Init() error {
	cfg := openai.DefaultConfig(config.GetConf().OpenAIConfig.Token)
	cfg.BaseURL = config.GetConf().OpenAIConfig.URL
	openAICli = openai.NewClientWithConfig(cfg)

	openAICache = cache.New(
		time.Duration(config.GetConf().OpenAICache.ExpireTime)*time.Second,
		time.Duration(config.GetConf().OpenAICache.CleanTime)*time.Second,
	)
	openAIWorkers = make(map[uint64]chan struct{})
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
	now := time.Now()
	ch := assignWork(msgID, reqMsg)
	select {
	case <-ctx.Done():
		logrus.Debugf("Talk:%v cost:%v, ctx.done:%v", reqMsg, time.Since(now), ctx.Err())
		return "", ctx.Err()
	case <-ch:
		cacheMsg, _ := getTalkCache(msgID)
		logrus.Debugf("Talk get cache:%v is:%v", msgID, cacheMsg)
		return cacheMsg, nil
	}
}

func assignWork(msgID uint64, msg string) chan struct{} {
	var ch chan struct{}
	var ok bool
	openAILock.Lock()
	ch, ok = openAIWorkers[msgID]
	openAILock.Unlock()
	if ok {
		return ch
	}

	// 加到 worker 里
	ch = make(chan struct{})
	openAILock.Lock()
	openAIWorkers[msgID] = ch
	openAILock.Unlock()

	go func() {
		// 关闭 ch 标识该 goroutine 的任务执行完了
		defer close(ch)
		now := time.Now()
		rsp, err := openAICli.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		})
		if err != nil {
			logrus.Errorf("Talk:%v err:%v, cost:%v", msg, err, time.Since(now))
			// 不用写缓存
			return
		}
		logrus.Debugf("Talk:%v:%v cost:%v, rsp:%v", msg, msgID, time.Since(now), rsp.Choices[0].Message.Content)
		// 写到缓存
		addTalkCache(msgID, rsp.Choices[0].Message.Content)
	}()
	return ch
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
	openAICache.Add(id, msg, time.Duration(cfg.OpenAICache.ExpireTime)*time.Second)
}
