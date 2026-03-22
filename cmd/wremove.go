package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/picker"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

func init() {
	wremoveCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return projectNameCompletion(cmd, args, toComplete)
		}
		return worktreeNameCompletion(cmd, args, toComplete)
	}
}

var wremoveCmd = &cobra.Command{
	Use:   "wremove [project-name] [worktree-name]",
	Short: "Remove a worktree from a project",
	Long: `Remove a worktree from a project. If only worktree name is provided, 
a project will be prompted. The main worktree cannot be removed.`,
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
			p, err = selectProject(cfg)
			if err != nil {
				return err
			}
			worktrees, err := worktree.ListWorktrees(cfg, p)
			if err != nil {
				return err
			}
			for _, wt := range worktrees {
				if wt.Name == args[0] {
					target = wt
					break
				}
			}
			if target == nil {
				return fmt.Errorf("worktree %q not found in project %q", args[0], p.Name)
			}
			if target.IsMain {
				return fmt.Errorf("cannot remove the main worktree of project %q", p.Name)
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

		question := fmt.Sprintf("%s Remove worktree %s at %s?",
			WarningStyle.Render("!"),
			ProjectNameStyle.Render(fmt.Sprintf("%q", target.Name)),
			PathStyle.Render(target.Path))

		confirmed, err := picker.Confirm(question)
		if err != nil {
			if errors.Is(err, picker.ErrConfirmationCancelled) {
				cmd.Println(MutedStyle.Render("Aborted."))
				return nil
			}
			return err
		}

		if !confirmed {
			cmd.Println(MutedStyle.Render("Aborted."))
			return nil
		}

		if err := worktree.RemoveWorktree(cfg, p, target.Name); err != nil {
			return err
		}

		cmd.Printf("%s %s %s %s\n",
			SuccessStyle.Render("Removed worktree"),
			ProjectNameStyle.Render(fmt.Sprintf("%q", target.Name)),
			MutedStyle.Render("from project"),
			ProjectNameStyle.Render(fmt.Sprintf("%q", p.Name)))
		return nil
	},
}
