// Package main provides the doopl CLI.
package main

import (
	"errors"
	"os"

	"github.com/arashackdev/doopl/internal/config"
	"github.com/urfave/cli/v2"
)

// getAuthKey retrieves the DeepL API auth key from flag, config file, or environment variable.
// It checks in order: --auth-key flag, ~/.doopl/config.toml, DEEPL_AUTH_KEY env var.
// Returns an error if no key is found in any source.
func getAuthKey(c *cli.Context) (string, error) {
	if key := c.String("auth-key"); key != "" {
		return key, nil
	}

	cfg, err := config.Load()
	if err != nil && !errors.Is(err, config.ErrNotFound) {
		return "", err
	}
	if cfg != nil && cfg.Auth.Key != "" {
		return cfg.Auth.Key, nil
	}

	if key := os.Getenv("DEEPL_AUTH_KEY"); key != "" {
		return key, nil
	}

	return "", errors.New("auth key not found: use doopl config set, --auth-key flag, or DEEPL_AUTH_KEY env var")
}
