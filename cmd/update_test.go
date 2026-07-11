package cmd

import (
	"testing"
)

func TestDetectPackageManager_Snap(t *testing.T) {
	msg, ok := detectPackageManager("/snap/sun/12/sun")
	if !ok {
		t.Fatal("expected snap detection")
	}
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
}

func TestDetectPackageManager_Homebrew(t *testing.T) {
	cases := []string{
		"/homebrew/Cellar/sun/0.3.0/bin/sun",
		"/home/linuxbrew/.linuxbrew/bin/sun",
	}
	for _, path := range cases {
		_, ok := detectPackageManager(path)
		if !ok {
			t.Errorf("expected homebrew detection for %s", path)
		}
	}
}

func TestDetectPackageManager_Nix(t *testing.T) {
	_, ok := detectPackageManager("/nix/store/abc-sun/bin/sun")
	if !ok {
		t.Fatal("expected nix detection")
	}
}

func TestDetectPackageManager_Normal(t *testing.T) {
	_, ok := detectPackageManager("/home/user/.local/bin/sun")
	if ok {
		t.Fatal("should not detect package manager for normal path")
	}
}

func TestDetectPackageManager_Linux(t *testing.T) {
	_, ok := detectPackageManager("/usr/local/bin/sun")
	if ok {
		t.Fatal("should not detect package manager for /usr/local/bin")
	}
}

func TestForceUpdateFlag_Registered(t *testing.T) {
	f := updateCmd.Flags().Lookup("force")
	if f == nil {
		t.Fatal("--force flag not registered on update command")
	}
	if f.DefValue != "false" {
		t.Fatalf("--force default = %q, want \"false\"", f.DefValue)
	}
}
