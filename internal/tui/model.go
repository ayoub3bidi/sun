package tui

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ayoub3bidi/sun/registry"
	sc "github.com/ayoub3bidi/solar-commons"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type step int

const (
	stepProjectType step = iota
	stepMigrationChoice
	stepProjectName
	stepAPIVersion
	stepAppTitle
	stepOIDC
	stepI18n
	stepSearch
	stepConfirm
	stepScaffolding
	stepDone
)

type projectKind int

const (
	kindAPI projectKind = iota
	kindDocs
)

type scaffoldResultMsg struct {
	result sc.ScaffoldResult
	err    error
}

type model struct {
	step           step
	projectKind    projectKind
	useAlembic     bool
	templateID     string
	projectName    textinput.Model
	apiVersion     textinput.Model
	appTitle       textinput.Model
	includeOIDC    bool
	oidcChosen     bool
	includeI18n    bool
	i18nChosen     bool
	includeSearch  bool
	searchChosen   bool
	validationErr  string
	spinner        spinner.Model
	err            error
	result         *sc.ScaffoldResult
	width          int
	height         int
}

func NewScaffoldModel() model {
	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Dot

	pn := textinput.New()
	pn.Placeholder = "my-project"
	pn.Focus()
	pn.CharLimit = 64
	pn.Width = 40

	av := textinput.New()
	av.Placeholder = "v1"
	av.CharLimit = 16
	av.Width = 40

	at := textinput.New()
	at.Placeholder = "My Project"
	at.CharLimit = 64
	at.Width = 40

	return model{
		step:        stepProjectType,
		projectName: pn,
		apiVersion:  av,
		appTitle:    at,
		spinner:     s,
	}
}

func Run(m model) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return err
	}
	if fm, ok := final.(model); ok && fm.err != nil {
		return fm.err
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		switch m.step {
		case stepProjectType:
			return m.updateProjectType(msg)
		case stepMigrationChoice:
			return m.updateMigrationChoice(msg)
		case stepProjectName:
			return m.updateProjectName(msg)
		case stepAPIVersion:
			return m.updateAPIVersion(msg)
		case stepAppTitle:
			return m.updateAppTitle(msg)
		case stepOIDC:
			return m.updateOIDC(msg)
		case stepI18n:
			return m.updateI18n(msg)
		case stepSearch:
			return m.updateSearch(msg)
		case stepConfirm:
			return m.updateConfirm(msg)
		case stepScaffolding:
			return m, nil
		case stepDone:
			if msg.String() == "enter" || msg.String() == "q" {
				return m, tea.Quit
			}
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case scaffoldResultMsg:
		if msg.err != nil {
			m.err = errors.New(friendlyScaffoldError(msg.err))
			m.step = stepDone
		} else {
			m.result = &msg.result
			m.step = stepDone
		}
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	m.validationErr = ""

	switch m.step {
	case stepProjectType:
		b.WriteString(titleStyle.Render("New Solar Project"))
		b.WriteString("\n\n")
		b.WriteString("What are you building?\n\n")
		b.WriteString(subtleStyle.Render("Press 'a' for API, 'd' for Docs"))

	case stepMigrationChoice:
		b.WriteString(titleStyle.Render("Database Migrations"))
		b.WriteString("\n\n")
		b.WriteString("Which migration style?\n\n")
		b.WriteString(subtleStyle.Render("'a' Alembic (recommended)  ·  'f' Flyway (classic SQL)"))

	case stepProjectName:
		b.WriteString(titleStyle.Render("Project Name"))
		b.WriteString("\n\n")
		b.WriteString(labelStyle.Render("Project name:"))
		b.WriteString("\n")
		b.WriteString(m.projectName.View())
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Kebab-case, e.g. my-api or my-docs"))

	case stepAPIVersion:
		b.WriteString(titleStyle.Render("API Version Prefix"))
		b.WriteString("\n\n")
		b.WriteString(labelStyle.Render("API version prefix:"))
		b.WriteString("\n")
		b.WriteString(m.apiVersion.View())
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("e.g. v1, v2, v2026-01"))

	case stepAppTitle:
		title := "Application Title"
		hint := "Used in page titles and API documentation"
		if m.projectKind == kindDocs {
			title = "Site Title"
			hint = "Shown in the navbar and browser tab"
		}
		b.WriteString(titleStyle.Render(title))
		b.WriteString("\n\n")
		b.WriteString(labelStyle.Render("Title:"))
		b.WriteString("\n")
		b.WriteString(m.appTitle.View())
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render(hint))

	case stepOIDC:
		b.WriteString(titleStyle.Render("OIDC Authentication"))
		b.WriteString("\n\n")
		b.WriteString("Include Google OIDC?")
		b.WriteString("\n\n")
		b.WriteString(boolChoice(m.includeOIDC, m.oidcChosen))
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Press 'y' for Yes, 'n' for No, then Enter"))

	case stepI18n:
		b.WriteString(titleStyle.Render("Internationalization"))
		b.WriteString("\n\n")
		b.WriteString("Include French (i18n)?")
		b.WriteString("\n\n")
		b.WriteString(boolChoice(m.includeI18n, m.i18nChosen))
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Press 'y' for Yes, 'n' for No, then Enter"))

	case stepSearch:
		b.WriteString(titleStyle.Render("Local Search"))
		b.WriteString("\n\n")
		b.WriteString("Include offline search plugin?")
		b.WriteString("\n\n")
		b.WriteString(boolChoice(m.includeSearch, m.searchChosen))
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Press 'y' for Yes, 'n' for No, then Enter"))

	case stepConfirm:
		tmpl, ok := m.selectedTemplate()
		if !ok {
			b.WriteString(errorStyle.Render("Template not found"))
			break
		}
		title := m.appTitle.Value()
		if title == "" {
			title = deriveTitle(m.projectName.Value())
		}
		b.WriteString(titleStyle.Render("Confirm Scaffold"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  Project name:  %s\n", m.projectName.Value()))
		if m.projectKind == kindAPI {
			b.WriteString(fmt.Sprintf("  API version:   %s\n", m.apiVersion.Value()))
			b.WriteString(fmt.Sprintf("  App title:     %s\n", title))
			b.WriteString(fmt.Sprintf("  OIDC:          %v\n", m.includeOIDC))
		} else {
			b.WriteString(fmt.Sprintf("  Site title:    %s\n", title))
			b.WriteString(fmt.Sprintf("  i18n:          %v\n", m.includeI18n))
			b.WriteString(fmt.Sprintf("  Search:        %v\n", m.includeSearch))
		}
		b.WriteString(fmt.Sprintf("  Template:      %s\n", tmpl.Label))
		b.WriteString("\n")
		b.WriteString(subtleStyle.Render("Press 'y' to scaffold, 'n' to go back"))

	case stepScaffolding:
		b.WriteString(m.spinner.View())
		b.WriteString(" Scaffolding your project...")

	case stepDone:
		if m.err != nil {
			b.WriteString(errorStyle.Render("Scaffold failed"))
			b.WriteString("\n\n")
			b.WriteString(m.err.Error())
			b.WriteString("\n\n")
			b.WriteString(subtleStyle.Render("Press Enter or 'q' to exit"))
		} else {
			b.WriteString(successStyle.Render("Project scaffolded!"))
			b.WriteString("\n\n")
			b.WriteString(fmt.Sprintf("  Location:  %s\n", m.result.Destination))
			b.WriteString(fmt.Sprintf("  Template:  %s\n", m.result.TemplateID))
			b.WriteString(fmt.Sprintf("  Commit:    %s\n", m.result.InitialCommit))
			if len(m.result.PrunedFeatures) > 0 {
				b.WriteString(fmt.Sprintf("  Pruned:    %s\n", strings.Join(m.result.PrunedFeatures, ", ")))
			}
			b.WriteString("\n\n")
			b.WriteString(subtleStyle.Render("Next steps:"))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf("  cd %s\n", m.projectName.Value()))
			if m.projectKind == kindAPI {
				b.WriteString("  cp .env.dist .env\n")
				b.WriteString("  docker compose up")
			} else {
				b.WriteString("  pnpm install\n")
				b.WriteString("  pnpm start")
			}
			b.WriteString("\n\n")
			b.WriteString(subtleStyle.Render("Press Enter or 'q' to exit"))
		}
	}

	if m.validationErr != "" {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(m.validationErr))
	}

	return docStyle.Render(b.String())
}

func boolChoice(yes, chosen bool) string {
	if yes {
		return "[✓] Yes"
	}
	if chosen {
		return "[✓] No"
	}
	return "[ ]"
}

func (m *model) selectedTemplate() (registry.Template, bool) {
	if m.projectKind == kindDocs {
		return registry.ByID("venus")
	}
	return registry.APIByMigration(m.useAlembic)
}

func (m model) updateProjectType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a", "A":
		m.projectKind = kindAPI
		m.step = stepMigrationChoice
	case "d", "D":
		m.projectKind = kindDocs
		m.templateID = "venus"
		m.includeI18n = true
		m.i18nChosen = true
		m.includeSearch = true
		m.searchChosen = true
		m.step = stepProjectName
		m.projectName.Focus()
	}
	return m, nil
}

func (m model) updateMigrationChoice(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a", "A":
		m.useAlembic = true
		m.templateID = "mars"
		m.step = stepProjectName
		m.projectName.Focus()
	case "f", "F":
		m.useAlembic = false
		m.templateID = "mercury"
		m.step = stepProjectName
		m.projectName.Focus()
	}
	return m, nil
}

func (m model) updateProjectName(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.projectName, cmd = m.projectName.Update(msg)
	if msg.String() == "enter" {
		name := strings.TrimSpace(m.projectName.Value())
		if err := validateProjectName(name); err != nil {
			m.validationErr = err.Error()
			return m, cmd
		}
		if destinationExists(name) {
			m.validationErr = fmt.Sprintf("directory %q already exists", name)
			return m, cmd
		}
		m.projectName.Blur()
		if m.projectKind == kindDocs {
			m.step = stepAppTitle
			m.appTitle.Focus()
			m.appTitle.SetValue(deriveTitle(name))
		} else {
			m.step = stepAPIVersion
			m.apiVersion.Focus()
			m.apiVersion.SetValue("v1")
		}
	}
	return m, cmd
}

func (m model) updateAPIVersion(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.apiVersion, cmd = m.apiVersion.Update(msg)
	if msg.String() == "enter" && strings.TrimSpace(m.apiVersion.Value()) != "" {
		m.apiVersion.Blur()
		m.step = stepAppTitle
		m.appTitle.Focus()
		m.appTitle.SetValue(deriveTitle(m.projectName.Value()))
	}
	return m, cmd
}

func (m model) updateAppTitle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.appTitle, cmd = m.appTitle.Update(msg)
	if msg.String() == "enter" {
		m.appTitle.Blur()
		if m.projectKind == kindAPI {
			m.step = stepOIDC
		} else {
			m.step = stepI18n
		}
	}
	return m, cmd
}

func (m model) updateOIDC(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.includeOIDC = true
		m.oidcChosen = true
	case "n", "N":
		m.includeOIDC = false
		m.oidcChosen = true
	case "enter":
		if m.oidcChosen {
			m.step = stepConfirm
		}
	}
	return m, nil
}

func (m model) updateI18n(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.includeI18n = true
		m.i18nChosen = true
	case "n", "N":
		m.includeI18n = false
		m.i18nChosen = true
	case "enter":
		if m.i18nChosen {
			m.step = stepSearch
		}
	}
	return m, nil
}

func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.includeSearch = true
		m.searchChosen = true
	case "n", "N":
		m.includeSearch = false
		m.searchChosen = true
	case "enter":
		if m.searchChosen {
			m.step = stepConfirm
		}
	}
	return m, nil
}

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.step = stepScaffolding
		return m, tea.Batch(m.spinner.Tick, scaffoldCmd(m))
	case "n", "N":
		m.step = stepProjectType
		m.oidcChosen = false
		m.i18nChosen = false
		m.searchChosen = false
		return m, nil
	}
	return m, nil
}

func scaffoldCmd(m model) tea.Cmd {
	return func() tea.Msg {
		tmpl, ok := m.selectedTemplate()
		if !ok {
			return scaffoldResultMsg{err: fmt.Errorf("template %q not found in registry", m.templateID)}
		}

		title := m.appTitle.Value()
		if title == "" {
			title = deriveTitle(m.projectName.Value())
		}
		dest := filepath.Join(".", m.projectName.Value())

		config := sc.ScaffoldConfig{
			Repo:          tmpl.Repo,
			Ref:           tmpl.DefaultRef,
			Destination:   dest,
			InitialCommit: "Initial scaffold",
		}

		if m.projectKind == kindAPI {
			config.Variables = map[string]string{
				"project_name": m.projectName.Value(),
				"api_version":  m.apiVersion.Value(),
				"app_title":    title,
			}
			if !m.includeOIDC {
				config.DisabledFeatures = append(config.DisabledFeatures, "oidc")
			}
		} else {
			config.Variables = map[string]string{
				"project_name": m.projectName.Value(),
				"site_title":   title,
			}
			if !m.includeI18n {
				config.DisabledFeatures = append(config.DisabledFeatures, "i18n")
			}
			if !m.includeSearch {
				config.DisabledFeatures = append(config.DisabledFeatures, "search")
			}
		}

		result, err := sc.Scaffold(config)
		return scaffoldResultMsg{result, err}
	}
}
