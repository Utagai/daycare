package cmd

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2/table"
	"github.com/spf13/cobra"
	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/project"
)

func timeAgo(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case d < 365*24*time.Hour:
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(d.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func init() {
	rootCmd.AddCommand(plistCmd)
}

var plistCmd = &cobra.Command{
	Use:   "plist",
	Short: "List all projects",
	Long:  "List all projects in the configured projects directory with their paths and last commit times.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		projects, err := project.ListProjects(cfg)
		if err != nil {
			return err
		}

		for _, p := range projects {
			if err := project.GetProjectInfo(p); err != nil {
				return fmt.Errorf("failed to get info for project %s: %w", p.Name, err)
			}
		}

		t := table.New().
			Headers("Name", "Path", "Last Commit")

		if len(projects) == 0 {
			cmd.Println("No projects found.")
			return nil
		}

		for _, p := range projects {
			t.Row(p.Name, p.Path, timeAgo(p.LastCommit))
		}

		cmd.Println(t.String())
		return nil
	},
}
