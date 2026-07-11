package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const checkInterval = 24 * time.Hour

type updateState struct {
	LastCheck time.Time `json:"last_check"`
}

var updateStatePath = func() string {
	if d, err := os.UserCacheDir(); err == nil {
		return filepath.Join(d, "sun", "update-check")
	}
	return filepath.Join(os.TempDir(), "sun-update-check")
}

func readUpdateState() updateState {
	var s updateState
	data, err := os.ReadFile(updateStatePath())
	if err != nil {
		return s
	}
	json.Unmarshal(data, &s)
	return s
}

func writeUpdateState(s updateState) {
	dir := filepath.Dir(updateStatePath())
	os.MkdirAll(dir, 0o755)
	data, _ := json.Marshal(s)
	os.WriteFile(updateStatePath(), data, 0o644)
}

func shouldCheck() bool {
	s := readUpdateState()
	return time.Since(s.LastCheck) >= checkInterval
}

func notifyIfUpdateAvailable() {
	if !shouldCheck() {
		return
	}
	writeUpdateState(updateState{LastCheck: time.Now()})

	latest, behind, err := checkForUpdate()
	if err != nil || !behind {
		return
	}

	fmt.Fprintf(os.Stderr, "\nsun %s is available (you have %s). Run 'sun update' to upgrade.\n\n", latest, normalizeVersion(version))
}
