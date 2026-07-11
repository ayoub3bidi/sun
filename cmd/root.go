package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = ""
	date    = ""
)

var checkUpdates bool

var rootCmd = &cobra.Command{
	Use:   "sun",
	Short: "Solar CLI — scaffold production-ready project templates",
	Long: `Sun is the Solar CLI for scaffolding project templates.

It supports Mars (FastAPI + Alembic), Mercury (FastAPI + Flyway),
and Venus (Docusaurus + Tailwind) templates.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if os.Getenv("SUN_CHECK_UPDATES") == "1" && cmd.Name() != "update" {
			go notifyIfUpdateAvailable()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sun %s\n", version)
		if commit != "" {
			fmt.Printf("commit: %s\n", commit)
		}
		if date != "" {
			fmt.Printf("built:  %s\n", date)
		}
		if checkUpdates {
			runVersionCheck()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	versionCmd.Flags().BoolVarP(&checkUpdates, "check", "c", false, "check for a newer version")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
}
