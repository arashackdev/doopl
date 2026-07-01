// Package main provides the doopl CLI.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	deepl "github.com/arashackdev/doopl/pkg/deepl"
	"github.com/arashackdev/doopl/pkg/model"
	"github.com/urfave/cli/v2"
)

// rephraseCommand returns the rephrase command for using the Write API.
func rephraseCommand() *cli.Command {
	return &cli.Command{
		Name:      "rephrase",
		Usage:     "rephrase text (Write API)",
		ArgsUsage: "TEXT...",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "lang", Required: true, Usage: "language code"},
			&cli.StringFlag{Name: "tone", Usage: "formal|informal|friendly", Value: "formal"},
			&cli.BoolFlag{Name: "emoji", Usage: "include emoji in output"},
		},
		Action: func(c *cli.Context) error {
			texts := c.Args().Slice()
			if len(texts) == 0 {
				return errors.New("no text provided")
			}

			authKey, err := getAuthKey(c)
			if err != nil {
				return err
			}

			client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-cli", deepl.Version))
			if err != nil {
				return err
			}

			var opts []deepl.RephraseOption
			if tone := c.String("tone"); tone != "" {
				opts = append(opts, deepl.WithTone(model.WriteTone(tone)))
			}
			if c.Bool("emoji") {
				opts = append(opts, deepl.WithEmojiMode(model.WriteEmojiAdd))
			}

			results, err := client.Rephrase(context.Background(), texts, c.String("lang"), opts...)
			if err != nil {
				return err
			}

			outfmt := c.String("output")
			switch outfmt {
			case "json":
				out, _ := json.MarshalIndent(results, "", "  ")
				fmt.Println(string(out))
			default:
				for _, r := range results {
					fmt.Println(r.Text)
				}
			}
			return nil
		},
	}
}
