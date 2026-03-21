package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

func init() {
	rootCmd.AddCommand(wremoveCmd)
}

var wremoveCmd = &cobra.Command{
	Use:   "wremove [project-name] [worktree-name]",
	Short: "Remove a worktree from a project",
	Long: `Remove a worktree from a project. If project name is not provided, 
it will be prompted. The main worktree cannot be removed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var p *models.Project
		var target *models.Worktree

		switch len(args) {
		case 0:
			// Use picker for both project and worktree
			p, err = selectProject(cfg)
			if err != nil {
				return err
			}
			// Exclude main worktree from selection
			target, err = selectWorktree(cfg, p, true)
			if err != nil {
				return err
			}
		case 1:
			// Project name provided, use picker for worktree
			p, err = project.GetProject(cfg, args[0])
			if err != nil {
				return err
			}
			// Exclude main worktree from selection
			target, err = selectWorktree(cfg, p, true)
			if err != nil {
				return err
			}
		case 2:
			p, err = project.GetProject(cfg, args[0])
			if err != nil {
				return err
			}

			worktreeName := args[1]
			worktrees, err := worktree.ListWorktrees(cfg, p)
			if err != nil {
				return err
			}

			for _, wt := range worktrees {
				if wt.Name == worktreeName {
					target = wt
					break
				}
			}

			if target == nil {
				return fmt.Errorf("worktree %q not found in project %q", worktreeName, p.Name)
			}

			if target.IsMain {
				return fmt.Errorf("cannot remove the main worktree of project %q", p.Name)
			}
		default:
			return fmt.Errorf("too many arguments")
		}

		cmd.Printf("Remove worktree %q at %s? (y/N): ", target.Name, target.Path)

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			cmd.Println("Aborted.")
			return nil
		}

		if err := worktree.RemoveWorktree(cfg, p, target.Name); err != nil {
			return err
		}

		cmd.Printf("Removed worktree %q from project %q.\n", target.Name, p.Name)
		return nil
	},
}
