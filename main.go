package main

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
)

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

	genkit.DefineFlow(
		"jokesFlow",
		func(ctx context.Context, theme string) (string, error) {
			gemini15pro := vertexai.Model("gemini-1.5-pro")
			//myJoke := core.LookupActionFor(atype.ActionType().Tool, "local", "myJoke")

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

			textResponse, _ := response.Text()
			return textResponse, nil
		},
	)

	//rsvp, _ := jokesFlow.Run(context.Background(), "abc")
	//fmt.Print(rsvp)
	if err := genkit.Init(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
}
