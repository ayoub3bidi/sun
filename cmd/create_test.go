package cmd

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	sc "github.com/ayoub3bidi/solar-commons"
)

func TestBuildScaffoldConfigDefaults(t *testing.T) {
	cmd := newCreateCmd()
	config, err := buildScaffoldConfig(createOptions{templateID: "mars", projectName: "my-api"}, cmd.Flags())
	if err != nil {
		t.Fatal(err)
	}
	if config.Destination != "my-api" {
		t.Errorf("Destination = %q, want %q", config.Destination, "my-api")
	}
	if config.Variables["app_title"] != "My Api" || config.Variables["api_version"] != "v1" {
		t.Errorf("unexpected Mars defaults: %#v", config.Variables)
	}
	if !reflect.DeepEqual(config.DisabledFeatures, []string{"oidc"}) {
		t.Errorf("DisabledFeatures = %#v, want [oidc]", config.DisabledFeatures)
	}
}

func TestBuildScaffoldConfigTemplateOptions(t *testing.T) {
	t.Run("Mercury includes OIDC when requested", func(t *testing.T) {
		cmd := newCreateCmd()
		if err := cmd.Flags().Set("with-oidc", "true"); err != nil {
			t.Fatal(err)
		}
		config, err := buildScaffoldConfig(createOptions{
			templateID: "mercury", projectName: "my-api", withOIDC: true,
		}, cmd.Flags())
		if err != nil {
			t.Fatal(err)
		}
		if len(config.DisabledFeatures) != 0 {
			t.Errorf("DisabledFeatures = %#v, want none", config.DisabledFeatures)
		}
	})

	t.Run("Venus prunes selected features", func(t *testing.T) {
		cmd := newCreateCmd()
		for _, flag := range []string{"without-i18n", "without-search"} {
			if err := cmd.Flags().Set(flag, "true"); err != nil {
				t.Fatal(err)
			}
		}
		config, err := buildScaffoldConfig(createOptions{
			templateID: "venus", projectName: "my-docs", withoutI18n: true, withoutSearch: true,
		}, cmd.Flags())
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(config.DisabledFeatures, []string{"i18n", "search"}) {
			t.Errorf("DisabledFeatures = %#v, want [i18n search]", config.DisabledFeatures)
		}
		if config.Variables["site_title"] != "My Docs" {
			t.Errorf("site_title = %q, want %q", config.Variables["site_title"], "My Docs")
		}
	})
}

func TestBuildScaffoldConfigRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name    string
		options createOptions
		flag    string
		wantErr string
	}{
		{name: "missing template", options: createOptions{projectName: "my-api"}, wantErr: "--template is required"},
		{name: "missing name", options: createOptions{templateID: "mars"}, wantErr: "--name is required"},
		{name: "unknown template", options: createOptions{templateID: "unknown", projectName: "my-api"}, wantErr: "unknown template"},
		{name: "invalid name", options: createOptions{templateID: "mars", projectName: "My API"}, wantErr: "invalid --name"},
		{name: "invalid format", options: createOptions{templateID: "mars", projectName: "my-api", format: "yaml"}, wantErr: "unsupported output format"},
		{name: "API flag on docs", options: createOptions{templateID: "venus", projectName: "my-docs"}, flag: "with-oidc", wantErr: "apply only to Mars and Mercury"},
		{name: "docs flag on API", options: createOptions{templateID: "mars", projectName: "my-api"}, flag: "without-search", wantErr: "apply only to the venus"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newCreateCmd()
			if tt.flag != "" {
				if err := cmd.Flags().Set(tt.flag, "true"); err != nil {
					t.Fatal(err)
				}
			}
			if tt.options.format == "" {
				tt.options.format = "text"
			}
			_, err := buildScaffoldConfig(tt.options, cmd.Flags())
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestWriteScaffoldResultJSON(t *testing.T) {
	var output bytes.Buffer
	err := writeScaffoldResult(&output, sc.ScaffoldResult{
		Destination:    "my-api",
		TemplateID:     "mars",
		Ref:            "v0.1.5",
		PrunedFeatures: []string{"oidc"},
		InitialCommit:  "Initial scaffold",
	}, "json")
	if err != nil {
		t.Fatal(err)
	}
	var got scaffoldOutput
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", output.String(), err)
	}
	if got.TemplateID != "mars" || got.Destination != "my-api" || got.Ref != "v0.1.5" || got.InitialCommit != "Initial scaffold" {
		t.Errorf("unexpected output: %#v", got)
	}
	if !reflect.DeepEqual(got.PrunedFeatures, []string{"oidc"}) {
		t.Errorf("PrunedFeatures = %#v, want [oidc]", got.PrunedFeatures)
	}
}

func TestWriteScaffoldResultJSONUsesEmptyFeatureArray(t *testing.T) {
	var output bytes.Buffer
	if err := writeScaffoldResult(&output, sc.ScaffoldResult{TemplateID: "mercury"}, "json"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"pruned_features":[]`) {
		t.Errorf("output = %q, want an empty pruned_features array", output.String())
	}
}

func TestWriteScaffoldResultText(t *testing.T) {
	var output bytes.Buffer
	if err := writeScaffoldResult(&output, sc.ScaffoldResult{TemplateID: "mars", Destination: "my-api"}, "text"); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "Scaffolded mars into my-api\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}
