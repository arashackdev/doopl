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

// Example: Document translation round-trip (upload, poll, download).
func exampleDocument() {
	client, err := doopl.New(os.Getenv("DEEPL_AUTH_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Open a document file
	file, err := os.Open("example.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Translate the document
	outputFile, err := os.Create("example_de.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = client.TranslateDocument(
		context.Background(),
		file,
		"example.txt",
		"DE", // target language
		outputFile,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Document translation complete!")
}
