// Package main provides the doopl CLI.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/arashackdev/doopl/pkg/model"
	"github.com/urfave/cli/v2"
)

// documentCommand returns the document command for managing document translation.
func documentCommand() *cli.Command {
	return &cli.Command{
		Name:  "document",
		Usage: "manage document translation",
		Subcommands: []*cli.Command{
			{
				Name:      "translate",
				Usage:     "translate a document (upload → poll → download)",
				ArgsUsage: "--source LANG --target LANG --input FILE --output FILE",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "source", Usage: "source language (auto-detect if omitted)"},
					&cli.StringFlag{Name: "target", Required: true, Usage: "target language"},
					&cli.StringFlag{Name: "input", Required: true, Usage: "input file path"},
					&cli.StringFlag{Name: "output", Required: true, Usage: "output file path"},
					&cli.StringFlag{Name: "glossary-id", Usage: "glossary ID"},
					&cli.StringFlag{Name: "style-id", Usage: "style preset ID"},
					&cli.StringFlag{Name: "translation-memory-id", Usage: "translation memory ID"},
					&cli.Float64Flag{Name: "translation-memory-threshold", Usage: "TM match threshold"},
				},
				Action: func(c *cli.Context) error {
					inputPath := c.String("input")
					outputPath := c.String("output")

					inputFile, err := os.Open(inputPath)
					defer inputFile.Close()
					if err != nil {
						return fmt.Errorf("cannot open input: %w", err)
					}

					outputFile, err := os.Create(outputPath)
					defer outputFile.Close()
					if err != nil {
						return fmt.Errorf("cannot create output: %w", err)
					}

					authKey, err := getAuthKey(c)
					if err != nil {
						return err
					}

					client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
					if err != nil {
						return err
					}

					var opts []deepl.DocumentUploadOption
					if v := c.String("source"); v != "" {
						opts = append(opts, deepl.WithDocumentSourceLang(v))
					}
					if v := c.String("glossary-id"); v != "" {
						opts = append(opts, deepl.WithDocumentGlossaryID(v))
					}
					if v := c.String("style-id"); v != "" {
						opts = append(opts, deepl.WithDocumentStyleID(v))
					}
					if v := c.String("translation-memory-id"); v != "" {
						opts = append(opts, deepl.WithDocumentTranslationMemoryID(v))
					}
					if c.IsSet("translation-memory-threshold") {
						threshold := c.Float64("translation-memory-threshold")
						opts = append(opts, deepl.WithDocumentTranslationMemoryThreshold(threshold))
					}

					fmt.Fprintf(os.Stderr, "Translating %s...\n", filepath.Base(inputPath))

					ctx := context.Background()
					err = client.TranslateDocument(ctx, inputFile, filepath.Base(inputPath), c.String("target"), outputFile, opts...)
					if err != nil {
						return err
					}

					fmt.Fprintf(os.Stderr, "✓ Saved to %s\n", outputPath)
					return nil
				},
			},
			{
				Name:      "upload",
				Usage:     "upload a document for translation (returns document ID and key)",
				ArgsUsage: "--source LANG --target LANG --input FILE",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "source", Usage: "source language (auto-detect if omitted)"},
					&cli.StringFlag{Name: "target", Required: true, Usage: "target language"},
					&cli.StringFlag{Name: "input", Required: true, Usage: "input file path"},
				},
				Action: func(c *cli.Context) error {
					inputPath := c.String("input")

					inputFile, err := os.Open(inputPath)
					defer inputFile.Close()
					if err != nil {
						return fmt.Errorf("cannot open input: %w", err)
					}

					authKey, err := getAuthKey(c)
					if err != nil {
						return err
					}

					client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
					if err != nil {
						return err
					}

					var opts []deepl.DocumentUploadOption
					if v := c.String("source"); v != "" {
						opts = append(opts, deepl.WithDocumentSourceLang(v))
					}

					handle, err := client.DocumentUpload(context.Background(), inputFile, filepath.Base(inputPath), c.String("target"), opts...)
					if err != nil {
						return err
					}

					out, _ := json.MarshalIndent(map[string]string{
						"document_id":  handle.DocumentID,
						"document_key": handle.DocumentKey,
					}, "", "  ")
					fmt.Println(string(out))
					return nil
				},
			},
			{
				Name:      "status",
				Usage:     "check document translation status",
				ArgsUsage: "--document-id ID --document-key KEY",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "document-id", Required: true, Usage: "document ID"},
					&cli.StringFlag{Name: "document-key", Required: true, Usage: "document key"},
				},
				Action: func(c *cli.Context) error {
					authKey, err := getAuthKey(c)
					if err != nil {
						return err
					}

					client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
					if err != nil {
						return err
					}

					handle := &model.DocumentHandle{
						DocumentID:  c.String("document-id"),
						DocumentKey: c.String("document-key"),
					}

					status, err := client.DocumentStatus(context.Background(), handle)
					if err != nil {
						return err
					}

					out, _ := json.MarshalIndent(status, "", "  ")
					fmt.Println(string(out))
					return nil
				},
			},
			{
				Name:      "download",
				Usage:     "download translated document",
				ArgsUsage: "--document-id ID --document-key KEY --output FILE",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "document-id", Required: true, Usage: "document ID"},
					&cli.StringFlag{Name: "document-key", Required: true, Usage: "document key"},
					&cli.StringFlag{Name: "output", Required: true, Usage: "output file path"},
				},
				Action: func(c *cli.Context) error {
					outputPath := c.String("output")

					outputFile, err := os.Create(outputPath)
					defer outputFile.Close()
					if err != nil {
						return fmt.Errorf("cannot create output: %w", err)
					}

					authKey, err := getAuthKey(c)
					if err != nil {
						return err
					}

					client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
					if err != nil {
						return err
					}

					handle := &model.DocumentHandle{
						DocumentID:  c.String("document-id"),
						DocumentKey: c.String("document-key"),
					}

					err = client.DocumentDownload(context.Background(), handle, outputFile)
					if err != nil {
						return err
					}

					fmt.Fprintf(os.Stderr, "✓ Saved to %s\n", outputPath)
					return nil
				},
			},
		},
	}
}
