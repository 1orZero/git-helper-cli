package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	API    APIConfig    `toml:"api"`
	Branch BranchConfig `toml:"branch"`
}

type APIConfig struct {
	APIEndpoint string `toml:"api_endpoint"`
	APISecret   string `toml:"api_secret"`
	Model       string `toml:"model"`
}

type BranchConfig struct {
	Pattern              string `toml:"pattern"`
	DescriptionFormat    string `toml:"description_format"`
	MaxDescriptionLength int    `toml:"max_description_length"`
	NumSuggestions       int    `toml:"num_suggestions"`
}

func LoadConfig(configPath string) (Config, error) {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return Config{}, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".config", "git-helper-cli", "config.toml")
	}

	var config Config
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	if err := ValidateConfig(&config); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func ValidateConfig(config *Config) error {
	if config.API.APISecret == "" {
		return fmt.Errorf("APISecret is required")
	}

	// Set default values for optional fields
	if config.API.APIEndpoint == "" {
		config.API.APIEndpoint = "https://api.openai.com/v1" // Set a default API endpoint
	}

	if config.API.Model == "" {
		config.API.Model = "gpt-4o-mini" // Set a default model
	}

	if config.Branch.Pattern == "" {
		config.Branch.Pattern = "${date}/feature/${description}"
	}

	if config.Branch.DescriptionFormat == "" {
		config.Branch.DescriptionFormat = "kebab-case"
	}

	if config.Branch.MaxDescriptionLength == 0 {
		config.Branch.MaxDescriptionLength = 50
	}

	if config.Branch.NumSuggestions == 0 {
		config.Branch.NumSuggestions = 10
	}

	return nil
}
