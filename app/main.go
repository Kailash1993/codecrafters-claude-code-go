package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	readToolSpecification := openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "Read",
		Description: openai.String("Read and return the contents of a file"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"file_path": map[string]any{
					"type":        "string",
					"description": "The path to the file to read",
				},
			},
			"required": []string{"file_path"},
		},
	})

	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(prompt),
				},
			},
		},
	}

	// break until for loop to process tool calls
	for {
		client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
		resp, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:    "anthropic/claude-haiku-4.5",
				Messages: messages,
				Tools:    []openai.ChatCompletionToolUnionParam{readToolSpecification},
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}

		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(resp.Choices[0].Message.Content),
				},
			},
		})

		toolCalls := resp.Choices[0].Message.ToolCalls
		if len(toolCalls) > 0 {
			for i := range toolCalls {
				toolCall := toolCalls[i]
				execResult, err := execute_tool(toolCall)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error executing tool: %v\n", err)
					os.Exit(1)
				}

				messages = append(messages, openai.ChatCompletionMessageParamUnion{
					OfTool: &openai.ChatCompletionToolMessageParam{
						Content: openai.ChatCompletionToolMessageParamContentUnion{
							OfString: openai.String(execResult),
						},
						ToolCallID: toolCall.ID,
					},
				})

			}
		} else {
			fmt.Print(resp.Choices[0].Message.Content)
			break
		}
	}
}

func execute_tool(tool_call openai.ChatCompletionMessageToolCallUnion) (string, error) {
	fmt.Println("Executing tool call:")
	fmt.Printf("  Name: %s\n", tool_call.Function.Name)
	fmt.Printf("  Arguments: %v\n", tool_call.Function.Arguments)

	if tool_call.Function.Name == "Read" {
		var args struct {
			FilePath string `json:"file_path"`
		}
		if err := json.Unmarshal([]byte(tool_call.Function.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}

		content, err := os.ReadFile(args.FilePath)
		if err != nil {
			return "", fmt.Errorf("error reading file '%s': %v", args.FilePath, err)
		}

		return string(content[:]), nil
	} else {
		return "", fmt.Errorf("unknown tool: %s", tool_call.Function.Name)
	}
}

// messages = [{ role: "user", content: prompt }]

// loop:
//     response = call_api(messages)
//     append response message to messages

//     if response has no tool_calls:
//         print response.content
//         exit

//     for each tool_call in response.tool_calls:
//         result = execute_tool(tool_call)
//         append {
//             role: "tool",
//             tool_call_id: tool_call.id,
//             content: result
//         } to messages
