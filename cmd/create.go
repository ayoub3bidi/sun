package cmd

import (
	"fmt"

	"github.com/ayoub3bidi/sun/internal/tui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Scaffold a new project from a Solar template",
	Long: `Launches an interactive TUI to collect project details and
scaffold a new project from a Solar template (Mars, Mercury, or Venus).`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := tui.Run(tui.NewScaffoldModel()); err != nil {
			return fmt.Errorf("create scaffold: %w", err)
		}
		return nil
	},
}
