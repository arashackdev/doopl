//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
	"github.com/arashackdev/doopl/model"
)

// Example: Glossary management (create, list, translate with glossary).
func exampleGlossary() {
	client, err := doopl.New(os.Getenv("DEEPL_AUTH_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a glossary
	entries := model.GlossaryEntries{
		"API":     "Schnittstelle",
		"library": "Bibliothek",
		"user":    "Benutzer",
	}

	glos, err := client.CreateGlossary(
		context.Background(),
		"tech-terms",
		"EN",
		"DE",
		entries,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created glossary: %s (ID: %s)\n", glos.Name, glos.GlossaryID)

	// List glossaries
	allGlos, err := client.ListGlossaries(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Glossaries: %d\n", len(allGlos))

	// Translate with glossary
	results, err := client.TranslateText(
		context.Background(),
		[]string{"The API is powerful!"},
		"DE",
		doopl.WithSourceLang("EN"),
		doopl.WithGlossaryID(glos.GlossaryID),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Translation: %s\n", results[0].Text)

	// Cleanup: Delete glossary
	err = client.DeleteGlossary(context.Background(), glos.GlossaryID)
	if err != nil {
		log.Fatal(err)
	}
}
