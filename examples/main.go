package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Main entry point for runnable examples.
// Run with: go run examples/main.go examples/*.go [example-name]
//
// Available examples:
//   - translate: Basic text translation
//   - glossary: Glossary management (create, list, use, delete)
//   - document: Document translation workflow (reference)
//
// Environment:
//
//	DEEPL_AUTH_KEY: Your DeepL API key (required)
func main() {
	flag.Parse()

	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		log.Fatal("DEEPL_AUTH_KEY environment variable not set")
	}

	exampleName := flag.Arg(0)
	if exampleName == "" {
		exampleName = "translate"
	}

	fmt.Printf("🚀 Running %q example...\n\n", exampleName)

	switch exampleName {
	case "translate":
		exampleTranslate()
	case "glossary":
		exampleGlossary()
	case "document":
		exampleDocument()
	default:
		fmt.Fprintf(os.Stderr, "Unknown example: %q\n", exampleName)
		fmt.Fprintf(os.Stderr, "Available: translate, glossary, document\n")
		os.Exit(1)
	}
}
