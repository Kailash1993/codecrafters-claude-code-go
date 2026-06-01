package main

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type LLMClient struct {
	Client openai.Client
}

func NewLLMClient(apiKey, baseUrl string) *LLMClient {
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}
	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	return &LLMClient{Client: client}
}

func (c *LLMClient) ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return c.Client.Chat.Completions.New(ctx, params)
}
