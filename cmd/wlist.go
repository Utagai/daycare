package cmd

import (
	"fmt"

	"charm.land/lipgloss/v2/table"
	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

func init() {
	rootCmd.AddCommand(wlistCmd)
}

var wlistCmd = &cobra.Command{
	Use:   "wlist [project-name]",
	Short: "List worktrees for a project",
	Long:  "List all worktrees for the specified project. If project name is not provided, select from available projects.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var p *models.Project
		if len(args) == 0 {
			return fmt.Errorf("automatic project selection not implemented yet. Please specify a project name.")
		}

		p, err = project.GetProject(cfg, args[0])
		if err != nil {
			return err
		}

		worktrees, err := worktree.ListWorktrees(cfg, p)
		if err != nil {
			return err
		}

		for _, wt := range worktrees {
			if err := worktree.GetWorktreeInfo(wt); err != nil {
				return fmt.Errorf("failed to get info for worktree %s: %w", wt.Name, err)
			}
		}

		t := table.New().
			Headers("Name", "Path", "Last Commit")

		if len(worktrees) == 0 {
			cmd.Println("No worktrees found.")
			return nil
		}

		for _, wt := range worktrees {
			t.Row(wt.Name, wt.Path, timeAgo(wt.LastCommit))
		}

		cmd.Println(t.String())
		return nil
	},
}
