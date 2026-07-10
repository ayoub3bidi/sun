package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func TestProjectTypeArrowEnterSelectsDocs(t *testing.T) {
	m := NewScaffoldModel()

	updated, _ := m.updateProjectType(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	if m.projectCursor != 1 {
		t.Fatalf("projectCursor=%d, want 1", m.projectCursor)
	}

	updated, _ = m.updateProjectType(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if m.step != stepProjectName {
		t.Fatalf("step=%v, want stepProjectName", m.step)
	}
	if m.projectKind != kindDocs || m.templateID != "venus" {
		t.Fatalf("selected kind/template = %v/%q, want docs/venus", m.projectKind, m.templateID)
	}
	if !m.includeI18n || !m.includeSearch {
		t.Fatal("docs defaults should include i18n and search")
	}
}

func TestMigrationArrowEnterSelectsFlyway(t *testing.T) {
	m := NewScaffoldModel()

	updated, _ := m.updateMigrationChoice(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	updated, _ = m.updateMigrationChoice(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if m.step != stepProjectName {
		t.Fatalf("step=%v, want stepProjectName", m.step)
	}
	if m.useAlembic || m.templateID != "mercury" {
		t.Fatalf("selected migration = alembic:%v template:%q, want flyway/mercury", m.useAlembic, m.templateID)
	}
}

func TestOIDCEnterAcceptsDefaultNo(t *testing.T) {
	m := NewScaffoldModel()
	m.step = stepOIDC

	updated, _ := m.updateOIDC(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if m.step != stepConfirm {
		t.Fatalf("step=%v, want stepConfirm", m.step)
	}
	if m.includeOIDC {
		t.Fatal("OIDC should default to false")
	}
}
