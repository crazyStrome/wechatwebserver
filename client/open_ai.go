package client

import (
	"context"
	"time"
	"wechatwebserver/config"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

var openAICli *openai.Client

func Init() error {
	cfg := openai.DefaultConfig(config.GetConf().OpenAIConfig.Token)
	cfg.BaseURL = config.GetConf().OpenAIConfig.URL
	openAICli = openai.NewClientWithConfig(cfg)
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

func Talk(ctx context.Context, reqMsg string) (string, error) {
	ch := make(chan string)
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
		logrus.Infof("Talk:%v cost:%v, err:%v", reqMsg, time.Since(now), err)
		if err != nil {
			logrus.Errorf("Talk:%v err:%v", reqMsg, err)
			return
		}
		ch <- rsp.Choices[0].Message.Content
	}()
	select {
	case <-ctx.Done():
		logrus.Infof("Talk:%v cost:%v, ctx.done", reqMsg, time.Since(now))
		return "", nil
	case data := <-ch:
		return data, nil
	}
}
