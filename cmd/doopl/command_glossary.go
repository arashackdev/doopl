// Package main provides the doopl CLI.
package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/arashackdev/doopl/pkg/model"
	"github.com/urfave/cli/v2"
)

// glossaryCommand returns the glossary command for managing glossaries.
func glossaryCommand() *cli.Command {
	return &cli.Command{
		Name:  "glossary",
		Usage: "manage glossaries",
		Subcommands: []*cli.Command{
			{
				Name:      "create",
				Usage:     "create a new glossary",
				ArgsUsage: "--name NAME --source LANG --target LANG --entries FILE",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "glossary name"},
					&cli.StringFlag{Name: "source", Required: true, Usage: "source language"},
					&cli.StringFlag{Name: "target", Required: true, Usage: "target language"},
					&cli.StringFlag{Name: "entries", Required: true, Usage: "CSV/TSV file with entries (source,target)"},
				},
				Action: func(c *cli.Context) error {
					entriesPath := c.String("entries")
					fileData, err := os.ReadFile(entriesPath)
					if err != nil {
						return fmt.Errorf("cannot read entries file: %w", err)
					}

					// Parse CSV/TSV
					entries := model.GlossaryEntries{}
					r := csv.NewReader(strings.NewReader(string(fileData)))
					r.Comma = '\t'
					rows, err := r.ReadAll()
					if err != nil {
						// Try comma
						r = csv.NewReader(strings.NewReader(string(fileData)))
						r.Comma = ','
						rows, err = r.ReadAll()
						if err != nil {
							return fmt.Errorf("cannot parse entries: %w", err)
						}
					}

					for _, row := range rows {
						if len(row) < 2 {
							continue
						}
						entries[strings.TrimSpace(row[0])] = strings.TrimSpace(row[1])
					}

					if len(entries) == 0 {
						return errors.New("no entries found in file")
					}

					authKey, err := getAuthKey(c)
					if err != nil {
						return err
					}

					client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
					if err != nil {
						return err
					}

					glossary, err := client.CreateGlossary(context.Background(), c.String("name"), c.String("source"), c.String("target"), entries)
					if err != nil {
						return err
					}

					out, _ := json.MarshalIndent(glossary, "", "  ")
					fmt.Println(string(out))
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list all glossaries",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Value: "table", Usage: "json|table"},
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

					glossaries, err := client.ListGlossaries(context.Background())
					if err != nil {
						return err
					}

					if c.String("format") == "json" {
						out, _ := json.MarshalIndent(glossaries, "", "  ")
						fmt.Println(string(out))
					} else {
						fmt.Printf("%-40s %-4s %-4s\n", "ID", "SRC", "TGT")
						fmt.Println(strings.Repeat("-", 50))
						for _, g := range glossaries {
							fmt.Printf("%-40s %-4s %-4s\n", g.GlossaryID, g.SourceLang, g.TargetLang)
						}
					}
					return nil
				},
			},
			{
				Name:      "get",
				Usage:     "get glossary details",
				ArgsUsage: "--id GLOSSARY_ID",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "id", Required: true, Usage: "glossary ID"},
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

					glossary, err := client.GetGlossary(context.Background(), c.String("id"))
					if err != nil {
						return err
					}

					out, _ := json.MarshalIndent(glossary, "", "  ")
					fmt.Println(string(out))
					return nil
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a glossary",
				ArgsUsage: "--id GLOSSARY_ID",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "id", Required: true, Usage: "glossary ID"},
					&cli.BoolFlag{Name: "force", Usage: "skip confirmation"},
				},
				Action: func(c *cli.Context) error {
					if !c.Bool("force") {
						fmt.Printf("Delete glossary %s? (yes/no): ", c.String("id"))
						var confirm string
						fmt.Scanln(&confirm)
						if confirm != "yes" {
							return errors.New("cancelled")
						}
					}

					authKey, err := getAuthKey(c)
					if err != nil {
						return err
					}

					client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
					if err != nil {
						return err
					}

					err = client.DeleteGlossary(context.Background(), c.String("id"))
					if err != nil {
						return err
					}

					fmt.Println("✓ Glossary deleted")
					return nil
				},
			},
			{
				Name:      "entries",
				Usage:     "get entries from a glossary",
				ArgsUsage: "--id GLOSSARY_ID",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "id", Required: true, Usage: "glossary ID"},
					&cli.StringFlag{Name: "format", Value: "json", Usage: "json|csv"},
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

					entries, err := client.GlossaryEntries(context.Background(), c.String("id"))
					if err != nil {
						return err
					}

					if c.String("format") == "json" {
						out, _ := json.MarshalIndent(entries, "", "  ")
						fmt.Println(string(out))
					} else {
						for src, tgt := range entries {
							fmt.Printf("%s,%s\n", src, tgt)
						}
					}
					return nil
				},
			},
		},
	}
}
