package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/picker"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

// selectProject displays a picker for project selection and returns the selected project.
// It pre-selects the project containing the current working directory, if any.
func selectProject(cfg *config.Config) (*models.Project, error) {
	projects, err := project.ListProjects(cfg)
	if err != nil {
		return nil, err
	}

	items := make([]picker.Item, len(projects))
	for i, p := range projects {
		items[i] = picker.ProjectItem{
			Name: p.Name,
			Path: p.Path,
		}
	}

	selectedIdx := -1
	if cwd, err := os.Getwd(); err == nil {
		if inferred, err := project.FindProjectByPath(cfg, cwd); err == nil && inferred != nil {
			for i, p := range projects {
				if p.Name == inferred.Name {
					selectedIdx = i
					break
				}
			}
		}
	}

	selected, err := picker.Run(items, "Select a project", selectedIdx)
	if err != nil {
		return nil, err
	}

	projectItem := selected.(picker.ProjectItem)
	return project.GetProject(cfg, projectItem.Name)
}

// selectWorktree displays a picker for worktree selection and returns the selected worktree.
// If excludeMain is true, the main worktree will be excluded from the list.
// It pre-selects the worktree containing the current working directory, if any.
func selectWorktree(cfg *config.Config, p *models.Project, excludeMain bool) (*models.Worktree, error) {
	worktrees, err := worktree.ListWorktrees(cfg, p)
	if err != nil {
		return nil, err
	}

	var items []picker.Item
	var filteredWorktrees []*models.Worktree
	for _, wt := range worktrees {
		if excludeMain && wt.IsMain {
			continue
		}
		items = append(items, picker.WorktreeItem{
			Name:   wt.Name,
			Branch: wt.Branch,
			Path:   wt.Path,
			IsMain: wt.IsMain,
		})
		filteredWorktrees = append(filteredWorktrees, wt)
	}

	selectedIdx := -1
	if cwd, err := os.Getwd(); err == nil {
		for i, wt := range filteredWorktrees {
			if wt.Path == cwd || strings.HasPrefix(cwd, wt.Path+string(filepath.Separator)) {
				selectedIdx = i
				break
			}
		}
	}

	selected, err := picker.Run(items, "Select a worktree", selectedIdx)
	if err != nil {
		return nil, err
	}

	worktreeItem := selected.(picker.WorktreeItem)
	for _, wt := range filteredWorktrees {
		if wt.Name == worktreeItem.Name {
			return wt, nil
		}
	}

	return nil, nil // should not happen
}
