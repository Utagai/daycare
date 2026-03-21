package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

func init() {
	rootCmd.AddCommand(wcdCmd)
}

var wcdCmd = &cobra.Command{
	Use:   "wcd [project-name] [worktree-name]",
	Short: "Output shell script to cd into a worktree",
	Long: `Output a shell script that cd's into the specified worktree.
Designed to be used with eval: eval "$(daycare wcd)"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var p *models.Project
		var wt *models.Worktree

		switch len(args) {
		case 0:
			// Use picker for both project and worktree
			p, err = selectProject(cfg)
			if err != nil {
				return err
			}
			wt, err = selectWorktree(cfg, p, false)
			if err != nil {
				return err
			}
		case 1:
			// Project name provided, use picker for worktree
			p, err = project.GetProject(cfg, args[0])
			if err != nil {
				return err
			}
			wt, err = selectWorktree(cfg, p, false)
			if err != nil {
				return err
			}
		case 2:
			p, err = project.GetProject(cfg, args[0])
			if err != nil {
				return err
			}

			worktrees, err := worktree.ListWorktrees(cfg, p)
			if err != nil {
				return err
			}

			for _, w := range worktrees {
				if w.Name == args[1] {
					wt = w
					break
				}
			}
			if wt == nil {
				return fmt.Errorf("worktree %q not found in project %q", args[1], p.Name)
			}
		default:
			return fmt.Errorf("too many arguments")
		}

		outputWcdScript(p, wt)
		return nil
	},
}

func outputWcdScript(p *models.Project, wt *models.Worktree) {
	worktreeDir := wt.Path
	pwcdPath := p.Path + "/.daycare/pwcd"

	if _, err := os.Stat(pwcdPath); err == nil {
		fmt.Printf("if [ -f %q ]; then\n  . %q\nfi\n", pwcdPath, pwcdPath)
	}

	fmt.Printf("cd %q\n", worktreeDir)
}
