// With options: translate with formality, source language, and other parameters.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
)

func main() {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		log.Fatal("DEEPL_AUTH_KEY environment variable not set")
	}

	client, err := deepl.New(authKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Translate with formality option (German supports formality)
	results, err := client.TranslateText(
		context.Background(),
		[]string{"Can you help me?"},
		"DE",
		deepl.WithSourceLang("EN"),
		deepl.WithFormality(deepl.FormalityMore), // Formal version
	)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	fmt.Println("Translation with more formal tone:")
	for _, result := range results {
		fmt.Printf("  %s\n", result.Text)
	}

	// Translate with less formality
	results, err = client.TranslateText(
		context.Background(),
		[]string{"Can you help me?"},
		"DE",
		deepl.WithSourceLang("EN"),
		deepl.WithFormality(deepl.FormalityLess), // Informal version
	)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	fmt.Println("\nTranslation with less formal tone:")
	for _, result := range results {
		fmt.Printf("  %s\n", result.Text)
	}
}
