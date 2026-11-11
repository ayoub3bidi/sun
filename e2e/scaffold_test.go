//go:build integration

package e2e_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
)

type cliResult struct {
	Destination    string   `json:"destination"`
	TemplateID     string   `json:"template_id"`
	Ref            string   `json:"ref"`
	PrunedFeatures []string `json:"pruned_features"`
	InitialCommit  string   `json:"initial_commit"`
}

func TestCreateCLIEndToEnd(t *testing.T) {
	binary := buildSunBinary(t)
	tests := []struct {
		name          string
		args          []string
		project       string
		template      string
		format        string
		textOutput    string
		pruned        []string
		presentPaths  []string
		absentPaths   []string
		contentFile   string
		contentAbsent []string
	}{
		{
			name:        "Mars defaults",
			args:        []string{"create", "--template", "mars", "--name", "my-api", "--format", "json"},
			project:     "my-api",
			template:    "mars",
			pruned:      []string{"oidc"},
			absentPaths: []string{"src/routes/oidc", "src/controllers/oidc", "src/utils/oidc.py", "solar.manifest.yaml"},
			contentFile: ".env.dist",
			contentAbsent: []string{
				"OIDC", "GOOGLE_", "mars",
			},
		},
		{
			name:         "Mercury with OIDC",
			args:         []string{"create", "--template", "mercury", "--name", "my-api", "--with-oidc", "--api-version", "v2", "--title", "My Billing API"},
			project:      "my-api",
			template:     "mercury",
			format:       "text",
			textOutput:   "Scaffolded mercury into my-api\n",
			presentPaths: []string{"src/routes/oidc"},
			absentPaths:  []string{"solar.manifest.yaml"},
		},
		{
			name:        "Venus without optional features",
			args:        []string{"create", "--template", "venus", "--name", "my-docs", "--without-i18n", "--without-search", "--format", "json"},
			project:     "my-docs",
			template:    "venus",
			pruned:      []string{"i18n", "search"},
			absentPaths: []string{"i18n/fr", "src/plugins/local-search.cjs", "solar.manifest.yaml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workDir := t.TempDir()
			stdout, stderr, err := runSun(binary, workDir, tt.args...)
			if err != nil {
				t.Fatalf("sun %s failed: %v\nstderr: %s", strings.Join(tt.args, " "), err, stderr)
			}
			if tt.format == "text" {
				if stdout != tt.textOutput {
					t.Errorf("stdout = %q, want %q", stdout, tt.textOutput)
				}
			} else {
				var result cliResult
				if err := json.Unmarshal([]byte(stdout), &result); err != nil {
					t.Fatalf("stdout is not valid JSON (%q): %v", stdout, err)
				}
				if result.TemplateID != tt.template {
					t.Errorf("template_id = %q, want %q", result.TemplateID, tt.template)
				}
				if result.Destination != tt.project {
					t.Errorf("destination = %q, want %q", result.Destination, tt.project)
				}
				if !equalStrings(result.PrunedFeatures, tt.pruned) {
					t.Errorf("pruned_features = %#v, want %#v", result.PrunedFeatures, tt.pruned)
				}
				if result.Ref == "" || result.InitialCommit != "Initial scaffold" {
					t.Errorf("unexpected result metadata: %#v", result)
				}
			}

			dest := filepath.Join(workDir, tt.project)
			for _, relative := range tt.presentPaths {
				if _, err := os.Stat(filepath.Join(dest, relative)); err != nil {
					t.Errorf("expected %s to exist: %v", relative, err)
				}
			}
			for _, relative := range tt.absentPaths {
				if _, err := os.Stat(filepath.Join(dest, relative)); !os.IsNotExist(err) {
					t.Errorf("expected %s to be absent, stat error: %v", relative, err)
				}
			}
			if tt.contentFile != "" {
				content, err := os.ReadFile(filepath.Join(dest, tt.contentFile))
				if err != nil {
					t.Fatal(err)
				}
				for _, unwanted := range tt.contentAbsent {
					if strings.Contains(string(content), unwanted) {
						t.Errorf("%s still contains %q", tt.contentFile, unwanted)
					}
				}
			}
			assertFreshGitRepository(t, dest)
		})
	}
}

func TestCreateCLIFailures(t *testing.T) {
	binary := buildSunBinary(t)
	tests := []struct {
		name      string
		args      []string
		project   string
		wantError string
	}{
		{name: "missing required name", args: []string{"create", "--template", "mars"}, wantError: "--name is required"},
		{name: "unknown template", args: []string{"create", "--template", "unknown", "--name", "my-app"}, project: "my-app", wantError: "unknown template"},
		{name: "feature flag not supported by template", args: []string{"create", "--template", "venus", "--name", "my-docs", "--with-oidc"}, project: "my-docs", wantError: "apply only to Mars and Mercury"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workDir := t.TempDir()
			_, stderr, err := runSun(binary, workDir, tt.args...)
			if err == nil {
				t.Fatal("expected command to fail")
			}
			if !strings.Contains(stderr, tt.wantError) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tt.wantError)
			}
			if tt.project != "" {
				if _, statErr := os.Stat(filepath.Join(workDir, tt.project)); !os.IsNotExist(statErr) {
					t.Errorf("failure left a scaffold directory, stat error: %v", statErr)
				}
			}
		})
	}

	t.Run("existing destination is preserved", func(t *testing.T) {
		workDir := t.TempDir()
		dest := filepath.Join(workDir, "existing")
		if err := os.Mkdir(dest, 0o755); err != nil {
			t.Fatal(err)
		}
		sentinel := filepath.Join(dest, "keep.txt")
		if err := os.WriteFile(sentinel, []byte("keep"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, stderr, err := runSun(binary, workDir, "create", "--template", "mars", "--name", "existing")
		if err == nil {
			t.Fatal("expected existing destination to fail")
		}
		if !strings.Contains(stderr, "destination already exists") {
			t.Errorf("stderr = %q, want existing destination error", stderr)
		}
		content, readErr := os.ReadFile(sentinel)
		if readErr != nil || string(content) != "keep" {
			t.Errorf("existing destination content changed: content=%q err=%v", content, readErr)
		}
	})
}

func buildSunBinary(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sun")
	cmd := exec.Command("go", "build", "-o", path, ".")
	cmd.Dir = ".."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build sun: %v\n%s", err, output)
	}
	return path
}

func runSun(binary, dir string, args ...string) (string, string, error) {
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func assertFreshGitRepository(t *testing.T, path string) {
	t.Helper()
	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatalf("open scaffold git repository: %v", err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatalf("read scaffold HEAD: %v", err)
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		t.Fatalf("read scaffold initial commit: %v", err)
	}
	if len(commit.ParentHashes) != 0 {
		t.Errorf("initial commit has %d parents, want none", len(commit.ParentHashes))
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
