package main

import "github.com/openai/openai-go/v3"

var readToolSpecification = openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
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
