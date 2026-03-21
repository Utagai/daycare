package cmd

import (
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/picker"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

// selectProject displays a picker for project selection and returns the selected project.
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

	selected, err := picker.Run(items, "Select a project")
	if err != nil {
		return nil, err
	}

	projectItem := selected.(picker.ProjectItem)
	return project.GetProject(cfg, projectItem.Name)
}

// selectWorktree displays a picker for worktree selection and returns the selected worktree.
// If excludeMain is true, the main worktree will be excluded from the list.
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

	selected, err := picker.Run(items, "Select a worktree")
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
