package main

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/ollama"
)

func main() {

	ctx := context.Background()

	if err := ollama.Init(ctx, "http://127.0.0.1:11434"); err != nil {
		log.Fatal(err)
	}

	gemmaModel := ollama.DefineModel(
		ollama.ModelDefinition{
			Name: "gemma2",
			Type: "generate",
		},
		&ai.ModelCapabilities{
			Multiturn:  false,
			SystemRole: true,
			Tools:      false,
			Media:      false,
		},
	)

	genRes, err := gemmaModel.Generate(ctx, ai.NewGenerateRequest(
		nil, ai.NewUserTextMessage("tell me a joke.")),
		func(ctx context.Context, grc *ai.GenerateResponseChunk) error {
			text, err := grc.Text()
			if err != nil {
				return err
			}
			fmt.Printf("%s", text)
			return nil
		})

	if err != nil {
		log.Fatal(err)
	}

	print("Response:")
	print(genRes.Text())
}
