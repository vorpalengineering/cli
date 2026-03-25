package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string `json:"api_key"`
	APIURL string `json:"api_url"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".vorpal")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func Path() string {
	return configPath()
}

func Load() (*Config, error) {
	cfg := &Config{
		APIURL: "http://localhost:8080",
	}

	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg, nil
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, nil
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir(), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}
