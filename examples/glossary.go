package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arashackdev/doopl/pkg/deepl"
	"github.com/arashackdev/doopl/pkg/model"
)

// Example: Glossary management (create, list, translate with glossary).
//
// Run with:
//   go run examples/glossary.go
//
// Environment:
//   DEEPL_AUTH_KEY: Your DeepL API key
//
// This example demonstrates:
// - Creating a new glossary with custom term pairs
// - Using a glossary in translation to enforce consistent terminology
// - Listing glossaries for your account
// - Deleting a glossary when done
func exampleGlossary() {
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

	ctx := context.Background()

	// Create a glossary with tech terms
	fmt.Println("📚 Creating glossary with tech terms...")
	glossary, err := client.CreateGlossary(
		ctx,
		"tech-terms",
		"EN",
		"DE",
		model.GlossaryEntries{
			"API":              "Schnittstelle",
			"frontend":         "Benutzeroberfläche",
			"backend":          "Backend-Server",
			"authentication":   "Authentifizierung",
			"encryption":       "Verschlüsselung",
		},
	)
	if err != nil {
		log.Fatal("create glossary:", err)
	}
	fmt.Printf("✓ Created glossary: %s\n\n", glossary.GlossaryID)

	// Use the glossary in translation
	fmt.Println("🌍 Translating with glossary...")
	results, err := client.TranslateText(
		ctx,
		[]string{
			"Our API provides a secure authentication system for the frontend and backend.",
		},
		"DE",
		deepl.WithSourceLang("EN"),
		deepl.WithGlossaryID(glossary.GlossaryID),
	)
	if err != nil {
		log.Fatal("translate:", err)
	}

	fmt.Printf("Source: %s\n", results[0].Text)
	fmt.Printf("German: %s\n\n", results[0].Text)

	// List glossaries
	fmt.Println("📋 Listing glossaries...")
	glossaries, err := client.ListGlossaries(ctx)
	if err != nil {
		log.Fatal("list glossaries:", err)
	}
	for _, g := range glossaries {
		fmt.Printf("  - %s (%s → %s)\n", g.Name, g.SourceLang, g.TargetLang)
	}

	// Clean up
	fmt.Printf("\n🗑️ Deleting glossary %s...\n", glossary.GlossaryID)
	err = client.DeleteGlossary(ctx, glossary.GlossaryID)
	if err != nil {
		log.Fatal("delete glossary:", err)
	}
	fmt.Println("✓ Done!\n")
}
