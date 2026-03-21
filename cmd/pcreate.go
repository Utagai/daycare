package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/project"
)

func init() {
	rootCmd.AddCommand(pcreateCmd)
}

var pcreateCmd = &cobra.Command{
	Use:   "pcreate [project-name]",
	Short: "Create a new project",
	Long: `Create a new project in the configured projects directory.
If project name is not provided, it will be prompted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter project name: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			name = strings.TrimSpace(input)
			if name == "" {
				return fmt.Errorf("project name cannot be empty")
			}
		}

		p, err := project.CreateProject(cfg, name)
		if err != nil {
			return err
		}

		cmd.Printf("%s %s at %s\n",
			SuccessStyle.Render("Created project"),
			ProjectNameStyle.Render(fmt.Sprintf("%q", p.Name)),
			PathStyle.Render(p.Path))
		return nil
	},
}
