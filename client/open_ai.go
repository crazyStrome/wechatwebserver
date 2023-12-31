package client

import (
	"context"
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
    rsp, err := openAICli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleUser,
                Content: reqMsg,
            },
        },
    })
    if err != nil {
        return "", err
    }
    return rsp.Choices[0].Message.Content, nil
}
