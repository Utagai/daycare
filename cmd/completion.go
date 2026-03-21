package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/project"
	"github.com/utagai/daycare/internal/worktree"
)

var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for zsh.

To load completions:

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it:

  autoload -U compinit
  compinit

  # To load completions for each session, execute once in your shell:
  # source <(daycare completion zsh)

  # To load completions for each session, once per login:
  # Add the above command to your ~/.zshrc
`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"zsh", "bash", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		default:
			return nil
		}
	},
}

func getProjectNames() ([]string, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	projects, err := project.ListProjects(cfg)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	return names, nil
}

func getWorktreeNames(projectName string) ([]string, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	proj, err := project.GetProject(cfg, projectName)
	if err != nil {
		return nil, err
	}

	worktrees, err := worktree.ListWorktrees(cfg, proj)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(worktrees))
	for i, w := range worktrees {
		names[i] = w.Name
	}
	return names, nil
}

func projectNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	names, err := getProjectNames()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func worktreeNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	names, err := getWorktreeNames(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
