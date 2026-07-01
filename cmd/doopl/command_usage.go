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

// usageCommand returns the usage command for checking API quota and usage.
func usageCommand() *cli.Command {
	return &cli.Command{
		Name:  "usage",
		Usage: "show API quota and usage",
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

			converter := &convert.ModelToEntityImpl{}
			row := converter.UsageRow(*usage)
			formatter := output.NewFormatter(c.String("output"))
			fmt.Print(formatter.FormatUsage(row))
			return nil
		},
	}
}
