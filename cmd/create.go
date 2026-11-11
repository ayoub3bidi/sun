package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	sc "github.com/ayoub3bidi/solar-commons"
	"github.com/ayoub3bidi/sun/internal/tui"
	"github.com/ayoub3bidi/sun/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createOptions struct {
	templateID    string
	projectName   string
	destination   string
	title         string
	apiVersion    string
	withOIDC      bool
	withoutI18n   bool
	withoutSearch bool
	format        string
}

type scaffoldOutput struct {
	Destination    string   `json:"destination"`
	TemplateID     string   `json:"template_id"`
	Ref            string   `json:"ref"`
	PrunedFeatures []string `json:"pruned_features"`
	InitialCommit  string   `json:"initial_commit"`
}

func newCreateCmd() *cobra.Command {
	var options createOptions
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "Scaffold a new project from a Solar template",
		Long:         "Scaffold a Solar project interactively, or provide flags to run non-interactively.",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				if err := tui.Run(tui.NewScaffoldModel()); err != nil {
					return fmt.Errorf("create scaffold: %w", err)
				}
				return nil
			}

			config, err := buildScaffoldConfig(options, cmd.Flags())
			if err != nil {
				return err
			}

			result, err := sc.Scaffold(config)
			if err != nil {
				if errors.Is(err, sc.ErrDestinationExists) {
					return fmt.Errorf("destination already exists: %s", config.Destination)
				}
				return fmt.Errorf("create scaffold: %w", err)
			}
			return writeScaffoldResult(cmd.OutOrStdout(), result, options.format)
		},
	}

	cmd.Flags().StringVar(&options.templateID, "template", "", "template ID (mars, mercury, or venus)")
	cmd.Flags().StringVar(&options.projectName, "name", "", "project name in lowercase kebab-case")
	cmd.Flags().StringVar(&options.destination, "destination", "", "output directory (defaults to ./<name>)")
	cmd.Flags().StringVar(&options.title, "title", "", "application or site title (defaults to a title derived from --name)")
	cmd.Flags().StringVar(&options.apiVersion, "api-version", "v1", "API version prefix for Mars and Mercury")
	cmd.Flags().BoolVar(&options.withOIDC, "with-oidc", false, "include Google OIDC support in Mars or Mercury")
	cmd.Flags().BoolVar(&options.withoutI18n, "without-i18n", false, "disable French i18n in Venus")
	cmd.Flags().BoolVar(&options.withoutSearch, "without-search", false, "disable offline search in Venus")
	cmd.Flags().StringVar(&options.format, "format", "text", "success output format: text or json")
	return cmd
}

func buildScaffoldConfig(options createOptions, flags *pflag.FlagSet) (sc.ScaffoldConfig, error) {
	options.projectName = strings.TrimSpace(options.projectName)
	if options.templateID == "" {
		return sc.ScaffoldConfig{}, errors.New("--template is required for non-interactive scaffolding")
	}
	if options.projectName == "" {
		return sc.ScaffoldConfig{}, errors.New("--name is required for non-interactive scaffolding")
	}
	if err := tui.ValidateProjectName(options.projectName); err != nil {
		return sc.ScaffoldConfig{}, fmt.Errorf("invalid --name: %w", err)
	}
	if options.format == "" {
		options.format = "text"
	}
	if options.format != "text" && options.format != "json" {
		return sc.ScaffoldConfig{}, fmt.Errorf("unsupported output format %q (use text or json)", options.format)
	}

	tmpl, ok := registry.ByID(options.templateID)
	if !ok {
		available := make([]string, 0, len(registry.List()))
		for _, candidate := range registry.List() {
			available = append(available, candidate.ID)
		}
		return sc.ScaffoldConfig{}, fmt.Errorf("unknown template %q (choose %s)", options.templateID, strings.Join(available, ", "))
	}

	apiTemplate := tmpl.Kind == "api"
	if apiTemplate {
		if flags.Changed("without-i18n") || flags.Changed("without-search") {
			return sc.ScaffoldConfig{}, fmt.Errorf("--without-i18n and --without-search apply only to the venus template")
		}
		if options.apiVersion == "" && flags.Changed("api-version") {
			return sc.ScaffoldConfig{}, errors.New("--api-version must not be empty")
		}
		if options.apiVersion == "" {
			options.apiVersion = "v1"
		}
	} else {
		if flags.Changed("with-oidc") || flags.Changed("api-version") {
			return sc.ScaffoldConfig{}, fmt.Errorf("--with-oidc and --api-version apply only to Mars and Mercury templates")
		}
	}

	destination := options.destination
	if destination == "" {
		destination = filepath.Join(".", options.projectName)
	}
	title := options.title
	if title == "" {
		title = tui.DeriveTitle(options.projectName)
	}

	config := sc.ScaffoldConfig{
		Repo:          tmpl.Repo,
		Ref:           tmpl.DefaultRef,
		Destination:   destination,
		InitialCommit: "Initial scaffold",
	}
	if apiTemplate {
		config.Variables = map[string]string{
			"project_name": options.projectName,
			"api_version":  options.apiVersion,
			"app_title":    title,
		}
		if !options.withOIDC {
			config.DisabledFeatures = []string{"oidc"}
		}
	} else {
		config.Variables = map[string]string{
			"project_name": options.projectName,
			"site_title":   title,
		}
		if options.withoutI18n {
			config.DisabledFeatures = append(config.DisabledFeatures, "i18n")
		}
		if options.withoutSearch {
			config.DisabledFeatures = append(config.DisabledFeatures, "search")
		}
	}
	return config, nil
}

func writeScaffoldResult(w io.Writer, result sc.ScaffoldResult, format string) error {
	if format == "json" {
		prunedFeatures := result.PrunedFeatures
		if prunedFeatures == nil {
			prunedFeatures = []string{}
		}
		return json.NewEncoder(w).Encode(scaffoldOutput{
			Destination:    result.Destination,
			TemplateID:     result.TemplateID,
			Ref:            result.Ref,
			PrunedFeatures: prunedFeatures,
			InitialCommit:  result.InitialCommit,
		})
	}
	_, err := fmt.Fprintf(w, "Scaffolded %s into %s\n", result.TemplateID, result.Destination)
	return err
}

var createCmd = newCreateCmd()
