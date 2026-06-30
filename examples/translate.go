//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
)

// Example: Basic text translation.
func main() {
	client, err := doopl.New(os.Getenv("DEEPL_AUTH_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	results, err := client.TranslateText(
		context.Background(),
		[]string{
			"Hello, world!",
			"Good morning!",
		},
		"DE", // target language
		doopl.WithFormality(doopl.FormalityMore),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i, r := range results {
		fmt.Printf("Result %d: %s\n", i, r.Text)
	}
}
