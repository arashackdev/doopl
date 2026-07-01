// Package main provides the doopl CLI.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/urfave/cli/v2"
)

// usageCommand returns the usage command for checking API quota and usage.
func usageCommand() *cli.Command {
	return &cli.Command{
		Name:  "usage",
		Usage: "show API quota and usage",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "format", Usage: "output format: json|table", Value: "table"},
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

			usage, err := client.Usage(context.Background())
			if err != nil {
				return err
			}

			if c.String("format") == "json" {
				out, _ := json.Marshal(usage)
				fmt.Println(string(out))
			} else {
				fmt.Printf("Characters:      %d / %d\n", usage.CharacterCount, usage.CharacterLimit)
				fmt.Printf("Documents:       %d / %d\n", usage.DocumentCount, usage.DocumentLimit)
				fmt.Printf("Team Documents:  %d / %d\n", usage.TeamDocumentCount, usage.TeamDocumentLimit)
			}
			return nil
		},
	}
}
