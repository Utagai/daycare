package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectsDir string `yaml:"projects_dir"`
}

// Load reads the configuration from ~/.config/daycare/config.yaml
// If the file doesn't exist, it returns a default configuration with ProjectsDir set to ~/projects
func Load() (*Config, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "daycare")
	configPath := filepath.Join(configDir, "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config
		return &Config{
			ProjectsDir: filepath.Join(os.Getenv("HOME"), "projects"),
		}, nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply default if projects_dir is not set
	if cfg.ProjectsDir == "" {
		cfg.ProjectsDir = filepath.Join(os.Getenv("HOME"), "projects")
	}

	// Expand ~ if present
	if cfg.ProjectsDir[0] == '~' {
		cfg.ProjectsDir = filepath.Join(os.Getenv("HOME"), cfg.ProjectsDir[1:])
	}

	return &cfg, nil
}

// Save writes the configuration to ~/.config/daycare/config.yaml
func (c *Config) Save() error {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "daycare")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
