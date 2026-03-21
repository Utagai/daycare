package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Print shell setup instructions",
	Long: `Print shell setup instructions to enable auto-cd functionality.

This command outputs shell code that can be added to your ~/.zshrc file.
It provides a wrapper function that intercepts 'pcd' and 'wcd' calls
and automatically evaluates their output to change directories.

Add the output of this command to your ~/.zshrc:

  echo "$(daycare setup)" >> ~/.zshrc

Then restart your shell or run 'source ~/.zshrc'.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		setupScript := `# Daycare shell integration - add to ~/.zshrc
daycare() {
  local prev=$1
  shift

  case "$prev" in
    pcd|wcd)
      local output
      output=$(command daycare "$@")
      if [ -n "$output" ]; then
        eval "$output"
      fi
      ;;
    *)
      command daycare "$@"
      ;;
  esac
}
`
		fmt.Println(setupScript)
		return nil
	},
}
