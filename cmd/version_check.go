package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	semver "github.com/Masterminds/semver/v3"
)

var githubLatestReleaseURL = "https://api.github.com/repos/ayoub3bidi/sun/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func normalizeVersion(v string) string {
	return strings.TrimPrefix(v, "v")
}

// isNewerVersion returns true when latestStr is strictly greater than currentStr
// using semver ordering. Falls back to string inequality when either value cannot
// be parsed (e.g. dev builds), so the behaviour is always defined.
func isNewerVersion(currentStr, latestStr string) bool {
	current, err1 := semver.NewVersion(currentStr)
	latest, err2 := semver.NewVersion(latestStr)
	if err1 != nil || err2 != nil {
		// Unparseable version — fall back to simple inequality so callers still
		// get a signal rather than silently swallowing the check.
		return latestStr != currentStr
	}
	return latest.GreaterThan(current)
}

func checkForUpdate() (latest string, behind bool, err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(githubLatestReleaseURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return "", false, fmt.Errorf("GitHub API rate limit exceeded")
	}
	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", false, fmt.Errorf("failed to parse response: %w", err)
	}

	latest = normalizeVersion(release.TagName)
	current := normalizeVersion(version)

	if latest == "" {
		return "", false, fmt.Errorf("empty version in latest release")
	}

	return latest, isNewerVersion(current, latest), nil
}

func runVersionCheck() {
	latest, behind, err := checkForUpdate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to check for updates: %v\n", err)
		return
	}
	if behind {
		fmt.Printf("update available: %s (you have %s). Run 'sun update' to upgrade.\n", latest, normalizeVersion(version))
	} else {
		fmt.Println("already up to date.")
	}
}
