// Package main provides the doopl CLI.
package main

import (
	"context"
	"fmt"

	"github.com/arashackdev/doopl/cmd/doopl/internal/convert"
	"github.com/arashackdev/doopl/cmd/doopl/internal/output"
	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/urfave/cli/v2"
)

// languagesCommand returns the languages command for listing supported languages.
func languagesCommand() *cli.Command {
	return &cli.Command{
		Name:      "languages",
		Usage:     "list supported languages",
		ArgsUsage: "[translate|document|glossary|write]",
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

			converter := &convert.ModelToEntityImpl{}
			rows := converter.LanguageRows(langs)
			formatter := output.NewFormatter(c.String("output"))
			fmt.Print(formatter.FormatLanguages(rows))
			return nil
		},
	}
}
