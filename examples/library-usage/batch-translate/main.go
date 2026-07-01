// Batch translate: translate multiple texts in a single request.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
	"github.com/arashackdev/doopl/pkg/model"
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

	// Translate multiple texts in one call (more efficient than individual requests)
	texts := []string{
		"Hello, world!",
		"How are you?",
		"Thank you for your help.",
	}

	results, err := client.TranslateText(
		context.Background(),
		texts,
		"DE", // German
	)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	// Print results (order preserved)
	fmt.Println("Batch translation results:")
	for i, result := range results {
		fmt.Printf("%d. %s => %s\n", i+1, texts[i], result.Text)
	}

	fmt.Printf("\nBilled characters: %d\n", sumBilled(results))
}

func sumBilled(results []model.TextResult) int {
	sum := 0
	for _, r := range results {
		sum += r.BilledCharacters
	}
	return sum
}
