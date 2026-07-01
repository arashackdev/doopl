package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
)

// Example: Basic text translation.
//
// Run with:
//
//	go run examples/translate.go
//
// Environment:
//
//	DEEPL_AUTH_KEY: Your DeepL API key
func exampleTranslate() {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		log.Fatal("DEEPL_AUTH_KEY environment variable not set")
	}

	client, err := deepl.New(authKey,
		deepl.WithAppInfo("doopl-example", "1.0.0"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Translate multiple texts with formality preference
	results, err := client.TranslateText(
		context.Background(),
		[]string{
			"Hello, world!",
			"Good morning!",
		},
		"DE", // German
		deepl.WithFormality(deepl.FormalityMore),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i, r := range results {
		fmt.Printf("[%d] %s\n", i, r.Text)
		fmt.Printf("    Source: %s, Billed: %d chars\n\n", r.DetectedSourceLang, r.BilledCharacters)
	}
}
