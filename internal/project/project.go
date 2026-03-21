package project

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
)

// isGitRepo returns true if the given directory is the root of a Git repository.
func isGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ListProjects scans the projects directory for Git repositories and returns them.
// The "worktrees" subdirectory is excluded from results.
func ListProjects(cfg *config.Config) ([]*models.Project, error) {
	entries, err := os.ReadDir(cfg.ProjectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("projects directory %q does not exist", cfg.ProjectsDir)
		}
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	var projects []*models.Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip the worktrees directory
		if entry.Name() == "worktrees" {
			continue
		}

		projectPath := filepath.Join(cfg.ProjectsDir, entry.Name())
		if !isGitRepo(projectPath) {
			continue
		}

		p := &models.Project{
			Name: entry.Name(),
			Path: projectPath,
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetProject returns a single project by name.
func GetProject(cfg *config.Config, name string) (*models.Project, error) {
	projectPath := filepath.Join(cfg.ProjectsDir, name)
	info, err := os.Stat(projectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project %q not found", name)
		}
		return nil, fmt.Errorf("failed to stat project directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", projectPath)
	}
	if !isGitRepo(projectPath) {
		return nil, fmt.Errorf("%q is not a Git repository", projectPath)
	}

	return &models.Project{
		Name: name,
		Path: projectPath,
	}, nil
}

// FindProjectByPath returns the project that contains the given path.
// It checks if the path is the project root or within its worktrees directory.
func FindProjectByPath(cfg *config.Config, targetPath string) (*models.Project, error) {
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return nil, err
	}
	absPath = filepath.Clean(absPath)

	projectsDir := filepath.Clean(cfg.ProjectsDir)

	// Check if target is within projectsDir
	if !strings.HasPrefix(absPath, projectsDir) {
		return nil, nil
	}

	relPath, err := filepath.Rel(projectsDir, absPath)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) < 1 {
		return nil, nil
	}

	projectName := parts[0]

	// Skip the worktrees directory itself
	if projectName == "worktrees" {
		if len(parts) < 2 {
			return nil, nil
		}
		projectName = parts[1]
	}

	// Verify the project exists
	projectPath := filepath.Join(projectsDir, projectName)
	if !isGitRepo(projectPath) {
		return nil, nil
	}

	return &models.Project{
		Name: projectName,
		Path: projectPath,
	}, nil
}

// GetProjectInfo populates the LastCommit field of a project by running `git log -1`.
func GetProjectInfo(p *models.Project) error {
	cmd := exec.Command("git", "log", "-1", "--format=%ct")
	cmd.Dir = p.Path

	out, err := cmd.Output()
	if err != nil {
		// Repo may have no commits yet; leave LastCommit as zero value.
		return nil
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil
	}

	var epoch int64
	if _, err := fmt.Sscanf(trimmed, "%d", &epoch); err != nil {
		return fmt.Errorf("failed to parse commit timestamp %q: %w", trimmed, err)
	}

	p.LastCommit = time.Unix(epoch, 0)
	return nil
}

// CreateProject creates a new project directory, initializes a Git repository,
// sets up the .daycare/ directory, and makes an initial commit.
func CreateProject(cfg *config.Config, name string) (*models.Project, error) {
	projectPath := filepath.Join(cfg.ProjectsDir, name)

	// Ensure projects directory exists
	if err := os.MkdirAll(cfg.ProjectsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create projects directory: %w", err)
	}

	// Create project directory
	if err := os.Mkdir(projectPath, 0755); err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("project %q already exists", name)
		}
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize Git repository
	if err := runGitCommand(projectPath, "git", "init"); err != nil {
		return nil, fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	// Create .daycare/ directory
	daycareDir := filepath.Join(projectPath, ".daycare")
	if err := os.Mkdir(daycareDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .daycare directory: %w", err)
	}

	// Create empty .daycare/pwcd script
	pwcdPath := filepath.Join(daycareDir, "pwcd")
	if err := createScript(pwcdPath); err != nil {
		return nil, fmt.Errorf("failed to create .daycare/pwcd: %w", err)
	}

	// Create empty .daycare/wcreate script
	wcreatePath := filepath.Join(daycareDir, "wcreate")
	if err := createScript(wcreatePath); err != nil {
		return nil, fmt.Errorf("failed to create .daycare/wcreate: %w", err)
	}

	// Add .daycare/ to .git/info/exclude for local ignore
	if err := addToGitExclude(projectPath, ".daycare/"); err != nil {
		return nil, fmt.Errorf("failed to add .daycare/ to git exclude: %w", err)
	}

	// Create initial commit
	if err := createInitialCommit(projectPath); err != nil {
		return nil, fmt.Errorf("failed to create initial commit: %w", err)
	}

	return &models.Project{
		Name: name,
		Path: projectPath,
	}, nil
}

// createScript creates an empty, executable shell script at the given path.
func createScript(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0755)
	if err != nil {
		return err
	}
	_, err = f.WriteString("#!/bin/sh\n")
	f.Close()
	return err
}

// addToGitExclude appends a pattern to .git/info/exclude.
func addToGitExclude(repoDir, pattern string) error {
	excludePath := filepath.Join(repoDir, ".git", "info", "exclude")

	// Ensure .git/info/ directory exists (it should after git init, but be safe)
	if err := os.MkdirAll(filepath.Dir(excludePath), 0755); err != nil {
		return err
	}

	// Check if pattern is already present
	data, err := os.ReadFile(excludePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == pattern {
			return nil // already present
		}
	}

	// Append the pattern
	f, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", pattern)
	return err
}

// createInitialCommit creates an initial empty commit in the repository.
func createInitialCommit(repoDir string) error {
	// Stage any existing files (there should be none, but be safe)
	if err := runGitCommand(repoDir, "git", "add", "--all"); err != nil {
		return err
	}

	// Allow empty commit in case there's nothing to stage
	if err := runGitCommand(repoDir, "git", "commit", "--allow-empty", "-m", "Initial commit"); err != nil {
		return err
	}

	return nil
}

// runGitCommand runs a git command in the given directory, returning an error
// that includes stderr output if the command fails.
func runGitCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("%w: %s", err, msg)
		}
		return err
	}
	return nil
}
