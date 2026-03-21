package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
	"github.com/utagai/daycare/internal/project"
)

func init() {
	rootCmd.AddCommand(pcdCmd)
}

var pcdCmd = &cobra.Command{
	Use:   "pcd [project-name]",
	Short: "Output shell script to cd into a project",
	Long: `Output a shell script that cd's into the specified project.
Designed to be used with eval: eval "$(daycare pcd)"`,
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

		outputPwcdScript(p)
		return nil
	},
}

func outputPwcdScript(p *models.Project) {
	projectDir := p.Path
	pwcdPath := p.Path + "/.daycare/pwcd"

	if _, err := os.Stat(pwcdPath); err == nil {
		fmt.Printf("if [ -f %q ]; then\n  . %q\nfi\n", pwcdPath, pwcdPath)
	}

	fmt.Printf("cd %q\n", projectDir)
}
