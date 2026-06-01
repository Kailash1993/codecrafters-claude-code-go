package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

type Executor struct {
	Client       *LLMClient
	ToolsHandler *ToolsHandler
}

func (e *Executor) Execute(model string, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam) (string, error) {
	resp, err := e.Client.ChatCompletion(context.Background(),
		openai.ChatCompletionNewParams{
			Model:    model,
			Messages: messages,
			Tools:    tools,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error calling ChatCompletion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	toolCalls := resp.Choices[0].Message.ToolCalls
	if len(toolCalls) > 0 {
		toolResults, err := e.ToolsHandler.HandleToolCalls(toolCalls)
		if err != nil {
			return "", fmt.Errorf("error handling tool calls: %w", err)
		}

		// For simplicity, we just return the first tool result. In a real implementation, you'd want to handle this more robustly.
		for _, result := range toolResults {
			return result, nil
		}
	}

	return resp.Choices[0].Message.Content, nil
}
