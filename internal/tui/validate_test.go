package tui

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "my-api", false},
		{"valid digits", "api2", false},
		{"empty", "", true},
		{"uppercase", "My-API", true},
		{"path sep", "../x", true},
		{"underscore", "my_api", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateProjectName(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("validateProjectName(%q) err=%v wantErr=%v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestDeriveTitle(t *testing.T) {
	if got := deriveTitle("my-api"); got != "My Api" {
		t.Fatalf("got %q", got)
	}
}
