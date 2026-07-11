package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"v0.3.0", "0.3.0"},
		{"0.3.0", "0.3.0"},
		{"v1.0.0-rc1", "1.0.0-rc1"},
	}
	for _, tc := range cases {
		if got := normalizeVersion(tc.input); got != tc.want {
			t.Errorf("normalizeVersion(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		current string
		latest  string
		want    bool
	}{
		// standard semver ordering
		{"0.3.0", "0.4.0", true},
		{"0.3.0", "0.3.0", false},
		{"0.4.0", "0.3.0", false},
		// multi-digit minor/patch — string comparison would get these wrong
		{"0.9.0", "0.10.0", true},
		{"0.9.0", "1.0.0", true},
		// pre-release: stable > pre-release of same version per semver
		{"1.0.0-rc1", "1.0.0", true},
		// pre-release behind another pre-release
		{"1.0.0-rc1", "1.0.0-rc2", true},
		// already on a newer pre-release than published stable
		{"1.1.0-dev", "1.0.0", false},
	}
	for _, tc := range cases {
		got := isNewerVersion(tc.current, tc.latest)
		if got != tc.want {
			t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tc.current, tc.latest, got, tc.want)
		}
	}
}

func TestCheckForUpdate_Behind(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(githubRelease{TagName: "v0.4.0"})
	}))
	defer srv.Close()

	orig := githubLatestReleaseURL
	githubLatestReleaseURL = srv.URL
	defer func() { githubLatestReleaseURL = orig }()

	version = "0.3.0"
	latest, behind, err := checkForUpdate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "0.4.0" {
		t.Fatalf("latest = %q, want 0.4.0", latest)
	}
	if !behind {
		t.Fatal("expected behind=true")
	}
}

func TestCheckForUpdate_UpToDate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(githubRelease{TagName: "v0.3.0"})
	}))
	defer srv.Close()

	orig := githubLatestReleaseURL
	githubLatestReleaseURL = srv.URL
	defer func() { githubLatestReleaseURL = orig }()

	version = "0.3.0"
	latest, behind, err := checkForUpdate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "0.3.0" {
		t.Fatalf("latest = %q, want 0.3.0", latest)
	}
	if behind {
		t.Fatal("expected behind=false")
	}
}

func TestCheckForUpdate_RateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	orig := githubLatestReleaseURL
	githubLatestReleaseURL = srv.URL
	defer func() { githubLatestReleaseURL = orig }()

	_, _, err := checkForUpdate()
	if err == nil {
		t.Fatal("expected error for rate limit")
	}
}

func TestCheckForUpdate_NetworkError(t *testing.T) {
	orig := githubLatestReleaseURL
	githubLatestReleaseURL = "http://127.0.0.1:1"
	defer func() { githubLatestReleaseURL = orig }()

	_, _, err := checkForUpdate()
	if err == nil {
		t.Fatal("expected error for network failure")
	}
}
