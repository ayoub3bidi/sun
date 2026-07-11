package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	selfupdate "github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var forceUpdate bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update sun to the latest version",
	Long: `Downloads and installs the latest release of sun from GitHub.

Replaces the current binary in-place. If sun was installed via a
package manager (Homebrew, Snap, Nix), use that manager to upgrade instead.`,
	SilenceUsage: true,
	RunE:         runUpdate,
}

func init() {
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "reinstall even if already on the latest version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate current binary: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("cannot resolve binary path: %w", err)
	}

	if msg, ok := detectPackageManager(exePath); ok {
		fmt.Println(msg)
		return nil
	}

	src, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return fmt.Errorf("failed to configure GitHub source: %w", err)
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source:  src,
		Filters: []string{fmt.Sprintf(`^sun_.*_%s_%s\.tar\.gz$`, runtime.GOOS, runtime.GOARCH)},
	})
	if err != nil {
		return fmt.Errorf("failed to configure updater: %w", err)
	}

	latest, found, err := updater.DetectLatest(context.Background(), selfupdate.ParseSlug("ayoub3bidi/sun"))
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	if !found {
		return fmt.Errorf("no release found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	current := normalizeVersion(version)
	latestVersion := normalizeVersion(latest.Version())

	if !forceUpdate && !isNewerVersion(current, latestVersion) {
		fmt.Printf("sun is already up to date (%s).\n", current)
		return nil
	}

	fmt.Printf("updating sun %s → %s ...\n", current, latestVersion)

	if err := updater.UpdateTo(context.Background(), latest, exePath); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied — try: sudo sun update\nor reinstall via: curl -fsSL https://raw.githubusercontent.com/ayoub3bidi/sun/main/scripts/install.sh | sh")
		}
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("updated sun %s → %s\n", current, latestVersion)
	return nil
}

func detectPackageManager(exePath string) (string, bool) {
	lower := strings.ToLower(exePath)
	switch {
	case strings.Contains(lower, "/snap/"):
		return "installed via Snap — run: sudo snap refresh sun", true
	case strings.Contains(lower, "/homebrew/") || strings.Contains(lower, "/cellar/") || strings.Contains(lower, "/linuxbrew/"):
		return "installed via Homebrew — run: brew upgrade sun", true
	case strings.Contains(lower, "/nix/"):
		return "installed via Nix — run: nix profile upgrade sun", true
	}
	return "", false
}
