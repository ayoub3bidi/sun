//go:build integration

package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	sc "github.com/ayoub3bidi/solar-commons"
	"github.com/ayoub3bidi/sun/registry"
	"github.com/go-git/go-git/v5"
)

func TestScaffoldMarsLikeTUI(t *testing.T) {
	tmpl, ok := registry.APIByMigration(true)
	if !ok {
		t.Fatal("mars template not in registry")
	}

	dest := filepath.Join(t.TempDir(), "my-api")
	result, err := sc.Scaffold(sc.ScaffoldConfig{
		Repo:        tmpl.Repo,
		Ref:         tmpl.DefaultRef,
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-api",
			"app_title":    "My API",
			"api_version":  "v1",
		},
		DisabledFeatures: []string{"oidc"},
		InitialCommit:    "Initial scaffold",
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if result.TemplateID != "mars" {
		t.Fatalf("TemplateID: got %q", result.TemplateID)
	}

	env, err := os.ReadFile(filepath.Join(dest, ".env.dist"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(env), "mars") {
		t.Error(".env.dist still contains mars")
	}
	if strings.Contains(string(env), "OIDC") || strings.Contains(string(env), "GOOGLE_") {
		t.Error(".env.dist still contains disabled OIDC settings")
	}
	for _, path := range []string{
		filepath.Join("src", "routes", "oidc"),
		filepath.Join("src", "controllers", "oidc"),
		filepath.Join("src", "utils", "oidc.py"),
		"solar.manifest.yaml",
	} {
		if _, err := os.Stat(filepath.Join(dest, path)); !os.IsNotExist(err) {
			t.Errorf("%s should have been removed", path)
		}
	}

	repo, err := git.PlainOpen(dest)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if len(commit.ParentHashes) != 0 {
		t.Error("expected single root commit")
	}
}

func TestScaffoldMercuryLikeTUI(t *testing.T) {
	tmpl, ok := registry.ByID("mercury")
	if !ok {
		t.Fatal("mercury template not in registry")
	}

	dest := filepath.Join(t.TempDir(), "my-flyway-api")
	result, err := sc.Scaffold(sc.ScaffoldConfig{
		Repo:        tmpl.Repo,
		Ref:         tmpl.DefaultRef,
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-flyway-api",
			"app_title":    "My Flyway API",
			"api_version":  "v1",
		},
		DisabledFeatures: []string{"oidc"},
		InitialCommit:    "Initial scaffold",
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if result.TemplateID != "mercury" {
		t.Fatalf("TemplateID: got %q", result.TemplateID)
	}

	env, err := os.ReadFile(filepath.Join(dest, ".env.dist"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(env), "mercury") || strings.Contains(string(env), "OIDC") || strings.Contains(string(env), "GOOGLE_") {
		t.Error("scaffold output still contains Mercury or disabled OIDC settings")
	}
	for _, path := range []string{
		filepath.Join("src", "routes", "oidc"),
		filepath.Join("src", "controllers", "oidc"),
		filepath.Join("src", "utils", "oidc.py"),
		"solar.manifest.yaml",
	} {
		if _, err := os.Stat(filepath.Join(dest, path)); !os.IsNotExist(err) {
			t.Errorf("%s should have been removed", path)
		}
	}

	repo, err := git.PlainOpen(dest)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if len(commit.ParentHashes) != 0 {
		t.Error("expected single root commit")
	}
}

func TestScaffoldVenusLikeTUI(t *testing.T) {
	tmpl, ok := registry.ByID("venus")
	if !ok {
		t.Fatal("venus template not in registry")
	}

	dest := filepath.Join(t.TempDir(), "my-docs")
	result, err := sc.Scaffold(sc.ScaffoldConfig{
		Repo:        tmpl.Repo,
		Ref:         tmpl.DefaultRef,
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-docs",
			"site_title":   "My Docs",
		},
		DisabledFeatures: []string{"i18n", "search"},
		InitialCommit:    "Initial scaffold",
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if result.TemplateID != "venus" {
		t.Fatalf("TemplateID: got %q", result.TemplateID)
	}

	for _, path := range []string{
		filepath.Join("i18n", "fr"),
		filepath.Join("src", "plugins", "local-search.cjs"),
		"solar.manifest.yaml",
	} {
		if _, err := os.Stat(filepath.Join(dest, path)); !os.IsNotExist(err) {
			t.Errorf("%s should have been removed", path)
		}
	}

	repo, err := git.PlainOpen(dest)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if len(commit.ParentHashes) != 0 {
		t.Error("expected single root commit")
	}
}

func TestRegistryTemplates(t *testing.T) {
	for _, id := range []string{"mars", "mercury", "venus"} {
		if _, ok := registry.ByID(id); !ok {
			t.Errorf("missing template %q", id)
		}
	}
}
