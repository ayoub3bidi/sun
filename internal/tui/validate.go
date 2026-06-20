package tui

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	sc "github.com/ayoub3bidi/solar-commons"
)

var projectNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func validateProjectName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("project name is required")
	}
	if strings.ContainsAny(name, `/\`) {
		return errors.New("project name must not contain path separators")
	}
	if !projectNamePattern.MatchString(name) {
		return errors.New("use kebab-case: lowercase letters, numbers, and hyphens (e.g. my-api)")
	}
	return nil
}

func destinationExists(name string) bool {
	dest := filepath.Join(".", name)
	info, err := os.Stat(dest)
	return err == nil && info.IsDir()
}

func friendlyScaffoldError(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, sc.ErrDestinationExists) {
		return "destination already exists — choose another project name or remove the existing directory"
	}
	return err.Error()
}

func deriveTitle(s string) string {
	if s == "" {
		return ""
	}
	words := strings.FieldsFunc(s, func(r rune) bool { return r == '-' || r == '_' })
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
