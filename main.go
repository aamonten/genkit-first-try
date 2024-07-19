package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/invopop/jsonschema"
)

type HelloPromptInput struct {
	UserName string
}

func main() {

	ctx := context.Background()

	if err := vertexai.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

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

	helloPrompt := ai.DefinePrompt(
		"prompts",
		"helloPrompt",
		nil,
		jsonschema.Reflect(HelloPromptInput{}),
		func(ctx context.Context, input any) (*ai.GenerateRequest, error) {
			params, ok := input.(HelloPromptInput)
			if !ok {
				return nil, errors.New("input doesn't satisfy schema")
			}
			prompt := fmt.Sprintf(
				"You are a helpful AI Assistant named Walt. Say hello %s, present yourself and tell a joke",
				params.UserName)
			return &ai.GenerateRequest{
				Messages: []*ai.Message{
					{Content: []*ai.Part{ai.NewTextPart(prompt)}},
				},
				Tools: []*ai.ToolDefinition{myJoke},
			}, nil
		},
	)

	genkit.DefineFlow(
		"jokesFlow",
		func(ctx context.Context, theme string) (string, error) {
			gemini15pro := vertexai.Model("gemini-1.5-pro")
			//myJoke := core.LookupActionFor(atype.ActionType().Tool, "local", "myJoke")

			request, err := helloPrompt.Render(context.Background(), HelloPromptInput{UserName: "Fred"})
			if err != nil {
				log.Fatal(err)
			}
			response, err := gemini15pro.Generate(ctx, request, nil)

			if err != nil {
				log.Fatal(err)
			}

			textResponse, _ := response.Text()
			return textResponse, nil
		},
	)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}
}
