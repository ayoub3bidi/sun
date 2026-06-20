package registry

import "testing"

func TestByID(t *testing.T) {
	tmpl, ok := ByID("mars")
	if !ok || tmpl.Kind != "api" {
		t.Fatal("expected mars api template")
	}
	if tmpl.Label == "" || tmpl.DefaultRef == "" {
		t.Fatal("mars template incomplete")
	}
}

func TestAPIByMigration(t *testing.T) {
	alembic, ok := APIByMigration(true)
	if !ok || alembic.ID != "mars" {
		t.Fatal("expected mars for alembic")
	}
	flyway, ok := APIByMigration(false)
	if !ok || flyway.ID != "mercury" {
		t.Fatal("expected mercury for flyway")
	}
}

func TestListHasThreeTemplates(t *testing.T) {
	if len(List()) != 3 {
		t.Fatalf("expected 3 templates, got %d", len(List()))
	}
}
