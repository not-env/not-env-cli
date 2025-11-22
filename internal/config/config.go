package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the CLI configuration
type Config struct {
	URL          string `toml:"url"`
	APIKey       string `toml:"api_key"`
	EnvID        *int64 `toml:"env_id,omitempty"`
	KeyType      string `toml:"key_type"`
	EnvIDFromKey *int64 `toml:"env_id_from_key"`
}

var configPath string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get home directory: %v", err))
	}
	configPath = filepath.Join(homeDir, ".not-env", "config")
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in. Run 'not-env login' first")
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Clear removes the configuration file
func Clear() error {
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config: %w", err)
	}
	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	return configPath
}

