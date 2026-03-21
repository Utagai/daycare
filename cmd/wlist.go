package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
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
			return fmt.Errorf("project name required: use 'wlist [project-name]' or use 'wselect' to pick a project")
		}

		p, err = project.GetProject(cfg, args[0])
		if err != nil {
			return err
		}

		worktrees, err := worktree.ListWorktrees(cfg, p)
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Name", "Branch", "Path"})
		t.SetStyle(table.StyleDefault)

		if len(worktrees) == 0 {
			cmd.Println("No worktrees found.")
			return nil
		}

		for _, wt := range worktrees {
			t.AppendRow(table.Row{wt.Name, wt.Branch, wt.Path})
		}

		cmd.Println(t.Render())
		return nil
	},
}
