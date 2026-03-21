package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set during build time
	Version = "0.1.0"
	// rootCmd is the main command
	rootCmd = &cobra.Command{
		Use:   "daycare",
		Short: "Daycare - Manage Git worktrees and projects",
		Long: `Daycare is a CLI tool for managing software projects and Git worktrees.
It helps you organize projects and quickly switch between different worktrees.`,
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(`{{.Name}} version {{.Version}}
`)

	rootCmd.AddCommand(pcreateCmd)
	rootCmd.AddCommand(plistCmd)
	rootCmd.AddCommand(pcdCmd)
	rootCmd.AddCommand(wcreateCmd)
	rootCmd.AddCommand(wlistCmd)
	rootCmd.AddCommand(wcdCmd)
	rootCmd.AddCommand(wremoveCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(setupCmd)
}

// Execute runs the root command
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
