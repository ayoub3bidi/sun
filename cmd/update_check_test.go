package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestShouldCheck_NoFile(t *testing.T) {
	dir := t.TempDir()
	orig := updateStatePath
	updateStatePath = func() string { return filepath.Join(dir, "update-check") }
	defer func() { updateStatePath = orig }()

	if !shouldCheck() {
		t.Fatal("should check when no state file exists")
	}
}

func TestShouldCheck_Fresh(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "update-check")
	os.WriteFile(path, []byte(`{"last_check":"`+time.Now().Format(time.RFC3339)+`"}`), 0o644)

	orig := updateStatePath
	updateStatePath = func() string { return path }
	defer func() { updateStatePath = orig }()

	if shouldCheck() {
		t.Fatal("should not check when state is fresh")
	}
}

func TestShouldCheck_Stale(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "update-check")
	stale := time.Now().Add(-25 * time.Hour)
	os.WriteFile(path, []byte(`{"last_check":"`+stale.Format(time.RFC3339)+`"}`), 0o644)

	orig := updateStatePath
	updateStatePath = func() string { return path }
	defer func() { updateStatePath = orig }()

	if !shouldCheck() {
		t.Fatal("should check when state is stale")
	}
}

func TestReadUpdateState_Empty(t *testing.T) {
	dir := t.TempDir()
	orig := updateStatePath
	updateStatePath = func() string { return filepath.Join(dir, "nonexistent") }
	defer func() { updateStatePath = orig }()

	s := readUpdateState()
	if !s.LastCheck.IsZero() {
		t.Fatal("expected zero time for missing file")
	}
}
