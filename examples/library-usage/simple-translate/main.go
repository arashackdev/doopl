// Simple example: translate a single text using doopl library.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
)

func main() {
	// Get API key from environment
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		log.Fatal("DEEPL_AUTH_KEY environment variable not set")
	}

	// Create a doopl client
	client, err := deepl.New(authKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Translate a single text
	results, err := client.TranslateText(context.Background(), []string{"Hello, world!"}, "DE")
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	// Print the result
	for _, result := range results {
		fmt.Printf("Original:  Hello, world!\n")
		fmt.Printf("Translated: %s\n", result.Text)
		fmt.Printf("Detected language: %s\n", result.DetectedSourceLang)
	}
}
