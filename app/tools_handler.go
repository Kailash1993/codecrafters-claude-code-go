package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go/v3"
)

type ToolsHandler struct {
}

func NewToolsHandler() *ToolsHandler {
	return &ToolsHandler{}
}

func (h *ToolsHandler) HandleToolCall(toolCall openai.ChatCompletionMessageToolCallUnion) (string, error) {
	switch toolCall.Function.Name {
	case "Read":
		var args struct {
			FilePath string `json:"file_path"`
		}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			log.Fatal("invalid arguments:", err)
		}

		// fmt.Printf("Reading file at path: %s\n", args.FilePath)
		content, err := os.ReadFile(args.FilePath)
		// fmt.Printf("File content read successfully, length: %d bytes\n", len(content))
		if err != nil {
			// fmt.Fprintf(os.Stderr, "error reading file '%s': %v\n", args.FilePath, err)
			os.Exit(1)
		}

		str1 := string(content[:])
		fmt.Print(str1)

		return str1, nil
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}
}

func (h *ToolsHandler) HandleToolCalls(toolCalls []openai.ChatCompletionMessageToolCallUnion) (map[string]string, error) {
	results := make(map[string]string)
	for _, toolCall := range toolCalls {
		result, err := h.HandleToolCall(toolCall)
		if err != nil {
			return nil, err
		}
		results[toolCall.Function.Name] = result
	}

	return results, nil
}
