package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	// Save and restore HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Create a temporary HOME directory
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Ensure no config file exists
	configPath := filepath.Join(tmpDir, ".config", "daycare", "config.yaml")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("config file should not exist: %v", err)
	}

	// Load config - should return defaults
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "code")
	if cfg.ProjectsDir != expectedPath {
		t.Errorf("ProjectsDir = %v, want %v", cfg.ProjectsDir, expectedPath)
	}
}

func TestLoadWithCustomConfig(t *testing.T) {
	// Save and restore HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Create a temporary HOME directory
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".config", "daycare")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `projects_dir: /custom/projects
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ProjectsDir != "/custom/projects" {
		t.Errorf("ProjectsDir = %v, want /custom/projects", cfg.ProjectsDir)
	}
}

func TestLoadWithEmptyProjectsDir(t *testing.T) {
	// Save and restore HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Create a temporary HOME directory
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Create config directory and file with empty projects_dir
	configDir := filepath.Join(tmpDir, ".config", "daycare")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `projects_dir: ""
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Load config - should use default
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "projects")
	if cfg.ProjectsDir != expectedPath {
		t.Errorf("ProjectsDir = %v, want %v", cfg.ProjectsDir, expectedPath)
	}
}

func TestLoadWithTildeExpansion(t *testing.T) {
	// Save and restore HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Create a temporary HOME directory
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Create config directory and file with tilde in path
	configDir := filepath.Join(tmpDir, ".config", "daycare")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `projects_dir: "~/myprojects"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Load config - tilde should be expanded
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "myprojects")
	if cfg.ProjectsDir != expectedPath {
		t.Errorf("ProjectsDir = %v, want %v", cfg.ProjectsDir, expectedPath)
	}
}

func TestSave(t *testing.T) {
	// Save and restore HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Create a temporary HOME directory
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Create config
	cfg := &Config{
		ProjectsDir: "/test/path",
	}

	// Save config
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tmpDir, ".config", "daycare", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	expected := "projects_dir: /test/path\n"
	if string(data) != expected {
		t.Errorf("config file content = %v, want %v", string(data), expected)
	}
}
