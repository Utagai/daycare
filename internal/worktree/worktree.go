package worktree

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
)

// worktreeDir returns the base directory where worktrees for a project are stored.
// Per the spec: <common_dir>/worktrees/<project-name>/
func worktreeDir(cfg *config.Config, projectName string) string {
	return filepath.Join(cfg.ProjectsDir, "worktrees", projectName)
}

// ListWorktrees returns all worktrees associated with the given project by
// parsing the output of `git worktree list --porcelain`.
func ListWorktrees(cfg *config.Config, p *models.Project) ([]*models.Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = p.Path

	out, err := cmd.Output()
	if err != nil {
		var stderr string
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = strings.TrimSpace(string(exitErr.Stderr))
		}
		if stderr != "" {
			return nil, fmt.Errorf("git worktree list failed: %s", stderr)
		}
		return nil, fmt.Errorf("git worktree list failed: %w", err)
	}

	return parseWorktreePorcelain(p.Name, out), nil
}

// parseWorktreePorcelain parses the `--porcelain` output of `git worktree list`.
// Each worktree block is separated by a blank line and has the form:
//
//	worktree <path>
//	HEAD <sha>
//	branch refs/heads/<branch>   (or "detached")
func parseWorktreePorcelain(projectName string, data []byte) []*models.Worktree {
	var worktrees []*models.Worktree

	scanner := bufio.NewScanner(bytes.NewReader(data))

	var (
		currentPath   string
		currentBranch string
		isFirst       = true
		firstPath     string
	)

	flush := func() {
		if currentPath == "" {
			return
		}
		name := filepath.Base(currentPath)
		isMain := currentPath == firstPath
		wt := &models.Worktree{
			Name:        name,
			Branch:      currentBranch,
			Path:        currentPath,
			ProjectName: projectName,
			IsMain:      isMain,
		}
		worktrees = append(worktrees, wt)
		currentPath = ""
		currentBranch = ""
	}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			flush()
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			currentPath = strings.TrimPrefix(line, "worktree ")
			if isFirst {
				firstPath = currentPath
				isFirst = false
			}
		} else if strings.HasPrefix(line, "branch ") {
			ref := strings.TrimPrefix(line, "branch ")
			// Convert refs/heads/<branch> -> <branch>
			currentBranch = strings.TrimPrefix(ref, "refs/heads/")
		} else if line == "detached" {
			currentBranch = "(detached)"
		}
	}

	// Flush the last block (no trailing blank line in some git versions)
	flush()

	return worktrees
}

// CreateWorktree creates a new Git worktree for the given project.
// The worktree is placed at <common_dir>/worktrees/<project-name>/<name>.
// A new branch named <name> is created for the worktree.
// Returns an error if a worktree with the given name already exists.
func CreateWorktree(cfg *config.Config, p *models.Project, name string) (*models.Worktree, error) {
	wtPath := filepath.Join(worktreeDir(cfg, p.Name), name)

	// Check that a worktree with this name doesn't already exist
	existing, err := ListWorktrees(cfg, p)
	if err != nil {
		return nil, err
	}
	for _, wt := range existing {
		if wt.Name == name {
			return nil, fmt.Errorf("worktree %q already exists for project %q", name, p.Name)
		}
	}

	// Ensure the parent worktrees directory exists
	if err := os.MkdirAll(worktreeDir(cfg, p.Name), 0755); err != nil {
		return nil, fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	// git worktree add <path> -b <branch>
	cmd := exec.Command("git", "worktree", "add", wtPath, "-b", name)
	cmd.Dir = p.Path

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("git worktree add failed: %s", msg)
		}
		return nil, fmt.Errorf("git worktree add failed: %w", err)
	}

	return &models.Worktree{
		Name:        name,
		Branch:      name,
		Path:        wtPath,
		ProjectName: p.Name,
		IsMain:      false,
	}, nil
}

// RemoveWorktree removes a Git worktree from the project.
// It refuses to remove the main worktree (the project directory itself).
func RemoveWorktree(cfg *config.Config, p *models.Project, name string) error {
	// Find the worktree to ensure it exists and isn't the main worktree
	worktrees, err := ListWorktrees(cfg, p)
	if err != nil {
		return err
	}

	var target *models.Worktree
	for _, wt := range worktrees {
		if wt.Name == name {
			target = wt
			break
		}
	}

	if target == nil {
		return fmt.Errorf("worktree %q not found in project %q", name, p.Name)
	}

	if target.IsMain {
		return fmt.Errorf("cannot remove the main worktree of project %q", p.Name)
	}

	// git worktree remove <path>
	cmd := exec.Command("git", "worktree", "remove", target.Path)
	cmd.Dir = p.Path

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("git worktree remove failed: %s", msg)
		}
		return fmt.Errorf("git worktree remove failed: %w", err)
	}

	return nil
}
