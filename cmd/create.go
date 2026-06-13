package cmd

import (
	"fmt"
	"os"

	"github.com/ayoub3bidi/sun/internal/tui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Scaffold a new project from a Solar template",
	Long: `Launches an interactive TUI to collect project details and
scaffold a new project from a Solar template (Mars, Mercury, or Venus).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		model := tui.NewScaffoldModel()
		if err := tui.Run(model); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return nil
	},
}
