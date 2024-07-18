package main

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/vertexai"
)

func main() {

	ctx := context.Background()

	if err := vertexai.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

	gemini15pro := vertexai.Model("gemini-1.5-pro")

	/* if err := ollama.Init(ctx, "http://127.0.0.1:11434"); err != nil {
		log.Fatal(err)
	} */

	/* 	gemmaModel := ollama.DefineModel(
		ollama.ModelDefinition{
			Name: "gemma2",
			Type: "generate",
		},
		&ai.ModelCapabilities{
			Multiturn:  false,
			SystemRole: true,
			Tools:      true,
			Media:      false,
		},
	) */

	myJoke := &ai.ToolDefinition{
		Name:        "myJoke",
		Description: "useful when you need a joke to tell",
		InputSchema: make(map[string]any),
		OutputSchema: map[string]any{
			"joke": "string",
		},
	}
	ai.DefineTool(
		myJoke,
		nil,
		func(ctx context.Context, input map[string]any) (map[string]any, error) {
			return map[string]any{"joke": "haha Just kidding no joke! got you"}, nil
		},
	)

	request := ai.GenerateRequest{
		Messages: []*ai.Message{
			{Content: []*ai.Part{ai.NewTextPart("Tell me a joke.")},
				Role: ai.RoleUser},
		},
		Tools: []*ai.ToolDefinition{myJoke},
	}
	response, err := gemini15pro.Generate(ctx, &request, nil)

	if err != nil {
		log.Fatal(err)
	}

	print("Response:")
	print(response.Text())
}
