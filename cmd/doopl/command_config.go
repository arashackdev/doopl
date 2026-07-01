// Package main provides the doopl CLI.
package main

import (
	"errors"
	"fmt"

	"github.com/arashackdev/doopl/internal/config"
	"github.com/urfave/cli/v2"
)

// configCommand returns the config command for managing doopl configuration.
func configCommand() *cli.Command {
	return &cli.Command{
		Name:    "config",
		Aliases: []string{"cfg"},
		Usage:   "manage doopl configuration",
		Subcommands: []*cli.Command{
			{
				Name:      "set",
				Usage:     "save auth key to ~/.doopl/config.toml",
				ArgsUsage: "AUTH_KEY [SERVER_URL]",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return errors.New("usage: doopl config set AUTH_KEY [SERVER_URL]")
					}
					cfg := &config.Config{}
					serverURL := ""
					if c.NArg() > 1 {
						serverURL = c.Args().Get(1)
					}
					return cfg.SetAuth(c.Args().Get(0), serverURL)
				},
			},
			{
				Name:  "show",
				Usage: "display current configuration",
				Action: func(_ *cli.Context) error {
					cfg, err := config.Load()
					if err != nil && !errors.Is(err, config.ErrNotFound) {
						return err
					}
					if cfg == nil {
						fmt.Println("no config found; use 'doopl config set' to initialize")
						return nil
					}
					if cfg.Auth.Key != "" {
						key := cfg.Auth.Key
						if len(key) > 6 {
							key = key[:3] + "..." + key[len(key)-3:]
						}
						fmt.Printf("Auth Key: %s\n", key)
					}
					if cfg.Auth.ServerURL != "" {
						fmt.Printf("Server URL: %s\n", cfg.Auth.ServerURL)
					}
					return nil
				},
			},
		},
	}
}
