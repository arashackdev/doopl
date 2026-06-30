// Package config handles persistent configuration for doopl CLI stored at ~/.doopl/config.toml.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the persistent doopl configuration.
type Config struct {
	Auth struct {
		Key       string `toml:"key"`
		ServerURL string `toml:"server_url,omitempty"`
	} `toml:"auth"`
}

var configDir = filepath.Join(os.Getenv("HOME"), ".doopl")

// ErrNotFound is returned when the config file doesn't exist.
var ErrNotFound = errors.New("config file not found")

// Load reads the config from ~/.doopl/config.toml, returning ErrNotFound if not present.
func Load() (*Config, error) {
	path := filepath.Join(configDir, "config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to ~/.doopl/config.toml, creating the directory if needed.
func (c *Config) Save() error {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	path := filepath.Join(configDir, "config.toml")
	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// AuthKey returns the auth key from config, env var DEEPL_AUTH_KEY, or fallback.
// Precedence: flag arg > env var > config file.
func (c *Config) AuthKey() string {
	if c != nil && c.Auth.Key != "" {
		return c.Auth.Key
	}
	return os.Getenv("DEEPL_AUTH_KEY")
}

// ServerURL returns the server URL from config, env var DEEPL_SERVER_URL, or empty string (use default).
func (c *Config) ServerURL() string {
	if c != nil && c.Auth.ServerURL != "" {
		return c.Auth.ServerURL
	}
	return os.Getenv("DEEPL_SERVER_URL")
}

// SetAuth updates the auth config and saves it.
func (c *Config) SetAuth(key string, serverURL string) error {
	c.Auth.Key = key
	c.Auth.ServerURL = serverURL
	return c.Save()
}
