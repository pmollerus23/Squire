package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ClientID  string `json:"clientId"`
	TenantID  string `json:"tenantId"`
	ServerURL string `json:"serverUrl"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".agent-cli-config.json")

	// Return defaults if config file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			ClientID:  "your-default-client-id",
			TenantID:  "your-default-tenant-id",
			ServerURL: "http://localhost:5000",
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".agent-cli-config.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
