package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var ErrNoAPIKey = errors.New("api_key is required in config")

type Config struct {
	APIKey string `toml:"api_key"`
	Model  string `toml:"model"`
}

func configDir() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine config directory: %w", err)
		}
		appData = filepath.Join(home, ".config")
	}
	return filepath.Join(appData, "prompt-enhancer"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	var cfg Config
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Model == "" {
		cfg.Model = "gemini-2.0-flash"
	}
	if cfg.APIKey == "" {
		return nil, ErrNoAPIKey
	}
	return &cfg, nil
}

func WriteConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	dir, _ := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}
