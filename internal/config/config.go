package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIBaseURL   string `json:"api_base_url"`
	AccessToken  string `json:"access_token"`
	Email        string `json:"email"`
}

var configDir = filepath.Join(os.Getenv("HOME"), ".groovekit")
var configFile = filepath.Join(configDir, "config.json")

// Load reads the config from ~/.groovekit/config.json
func Load() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist yet, return default
			return &Config{
				APIBaseURL: getAPIBaseURL(),
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Environment variable takes precedence
	if envURL := os.Getenv("GROOVEKIT_API_URL"); envURL != "" {
		cfg.APIBaseURL = envURL
	} else if cfg.APIBaseURL == "" {
		// Set default API URL if not present
		cfg.APIBaseURL = getAPIBaseURL()
	}

	return &cfg, nil
}

// getAPIBaseURL returns the API base URL from env var or default
func getAPIBaseURL() string {
	if envURL := os.Getenv("GROOVEKIT_API_URL"); envURL != "" {
		return envURL
	}
	return "https://api.groovekit.com"
}

// Save writes the config to ~/.groovekit/config.json
func (c *Config) Save() error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Clear removes the config file
func Clear() error {
	return os.Remove(configFile)
}

// IsAuthenticated checks if user is logged in
func (c *Config) IsAuthenticated() bool {
	return c.AccessToken != ""
}
