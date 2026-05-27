package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

func init() {
	wcreateCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return projectNameCompletion(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

var wcreateCmd = &cobra.Command{
	Use:   "wcreate [project-name] [name]",
	Short: "Create a new worktree",
	Long: `Create a new worktree for the specified project.
If project name is not provided, it will be prompted.
If worktree name is not provided, it will be prompted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var p *models.Project
		var worktreeName string

		switch len(args) {
		case 0:
			p, err = selectProject(cfg)
			if err != nil {
				return err
			}
		case 1:
			p, err = project.GetProject(cfg, args[0])
			if err != nil {
				return err
			}
		default:
			p, err = project.GetProject(cfg, args[0])
			if err != nil {
				return err
			}
			worktreeName = args[1]
		}

		if worktreeName == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter worktree name: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			worktreeName = strings.TrimSpace(input)
			if worktreeName == "" {
				return fmt.Errorf("worktree name cannot be empty")
			}
		}

		wt, err := worktree.CreateWorktree(cfg, p, worktreeName)
		if err != nil {
			return err
		}

		if err := runWcreateScript(p, wt); err != nil {
			return fmt.Errorf("failed to run wcreate script: %w", err)
		}

		cmd.Printf("%s %s %s %s\n",
			SuccessStyle.Render("Created worktree"),
			ProjectNameStyle.Render(fmt.Sprintf("%q", wt.Name)),
			MutedStyle.Render("for project"),
			ProjectNameStyle.Render(fmt.Sprintf("%q", p.Name)))
		return wcdCmd.RunE(cmd, []string{p.Name, wt.Name})
	},
}

func runWcreateScript(p *models.Project, wt *models.Worktree) error {
	scriptPath := filepath.Join(p.Path, ".daycare", "wcreate")
	info, err := os.Stat(scriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to stat wcreate script: %w", err)
	}
	if info.IsDir() {
		return nil
	}

	scriptPath, err = filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for wcreate script: %w", err)
	}

	c := exec.Command(scriptPath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	env := os.Environ()
	env = append(env, "PROJECT_NAME="+p.Name)
	env = append(env, "PROJECT_DIR="+p.Path)
	env = append(env, "WORKTREE_DIR="+wt.Path)
	c.Env = env

	if err := c.Run(); err != nil {
		if strings.Contains(err.Error(), "exit status 127") {
			return nil
		}
		return fmt.Errorf("wcreate script failed: %w", err)
	}

	return nil
}
