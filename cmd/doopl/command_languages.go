// Package main provides the doopl CLI.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/urfave/cli/v2"
)

// languagesCommand returns the languages command for listing supported languages.
func languagesCommand() *cli.Command {
	return &cli.Command{
		Name:      "languages",
		Usage:     "list supported languages",
		ArgsUsage: "[translate|document|glossary|write]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "format", Usage: "output format: json|table", Value: "table"},
		},
		Action: func(c *cli.Context) error {
			resource := "translate"
			if c.NArg() > 0 {
				resource = c.Args().Get(0)
			}

			authKey, err := getAuthKey(c)
			if err != nil {
				return err
			}

			client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
			if err != nil {
				return err
			}

			langs, err := client.Languages(context.Background(), resource)
			if err != nil {
				return err
			}

			if c.String("format") == "json" {
				out, _ := json.Marshal(langs)
				fmt.Println(string(out))
			} else {
				fmt.Printf("%-5s %s\n", "CODE", "NAME")
				fmt.Println(strings.Repeat("-", 50))
				for _, lang := range langs {
					fmt.Printf("%-5s %s\n", lang.Code, lang.Name)
				}
			}
			return nil
		},
	}
}
