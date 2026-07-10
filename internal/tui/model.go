package tui

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	sc "github.com/ayoub3bidi/solar-commons"
	"github.com/ayoub3bidi/sun/registry"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	step            step
	projectKind     projectKind
	useAlembic      bool
	templateID      string
	projectCursor   int
	migrationCursor int
	projectName     textinput.Model
	apiVersion      textinput.Model
	appTitle        textinput.Model
	includeOIDC     bool
	oidcChosen      bool
	includeI18n     bool
	i18nChosen      bool
	includeSearch   bool
	searchChosen    bool
	validationErr   string
	spinner         spinner.Model
	err             error
	result          *sc.ScaffoldResult
	width           int
	height          int
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
		inputWidth := m.inputWidth()
		m.projectName.Width = inputWidth
		m.apiVersion.Width = inputWidth
		m.appTitle.Width = inputWidth
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
	title, subtitle, body, footer := m.viewParts()
	if m.validationErr != "" {
		body += "\n\n" + errorStyle.Render(m.validationErr)
	}
	return m.renderFrame(title, subtitle, body, footer)
}

func (m model) viewParts() (string, string, string, string) {
	switch m.step {
	case stepProjectType:
		return "New Solar Project",
			"Choose the kind of Solar template to scaffold.",
			renderChoiceCards(projectTypeOptions(), m.projectCursor),
			"↑/↓ Move • Enter Select • A API • D Docs • ESC Quit"

	case stepMigrationChoice:
		return "Database Migrations",
			"Pick the migration workflow for the generated FastAPI backend.",
			renderChoiceCards(migrationOptions(), m.migrationCursor),
			"↑/↓ Move • Enter Select • A Alembic • F Flyway • ESC Quit"

	case stepProjectName:
		body := formField("Project name", m.projectName.View(), "Kebab-case only, for example my-api or my-docs.")
		return "Project Name", "This becomes the output directory and project identifier.", body, "Enter Continue • ESC Quit"

	case stepAPIVersion:
		body := formField("API version prefix", m.apiVersion.View(), "Examples: v1, v2, v2026-01.")
		return "API Version Prefix", "Used by the generated API routing and docs.", body, "Enter Continue • ESC Quit"

	case stepAppTitle:
		title := "Application Title"
		hint := "Used in page titles and API documentation"
		if m.projectKind == kindDocs {
			title = "Site Title"
			hint = "Shown in the navbar and browser tab"
		}
		body := formField("Title", m.appTitle.View(), hint+". Leave as generated or edit it.")
		return title, "Give the generated project a human-readable name.", body, "Enter Continue • ESC Quit"

	case stepOIDC:
		body := renderBooleanChoice("Include Google OIDC?", "Adds Google OIDC routes and configuration to the generated API.", m.includeOIDC)
		return "OIDC Authentication", "Choose whether this API starts with Google OIDC support.", body, "←/→ Move • Y/N Choose • Enter Continue • ESC Quit"

	case stepI18n:
		body := renderBooleanChoice("Include French i18n?", "Keeps the French locale and i18n project structure in the docs site.", m.includeI18n)
		return "Internationalization", "Choose whether the docs scaffold should include localization.", body, "←/→ Move • Y/N Choose • Enter Continue • ESC Quit"

	case stepSearch:
		body := renderBooleanChoice("Include offline search?", "Keeps the local search plugin in the generated docs site.", m.includeSearch)
		return "Local Search", "Choose whether readers should get built-in local search.", body, "←/→ Move • Y/N Choose • Enter Continue • ESC Quit"

	case stepConfirm:
		tmpl, ok := m.selectedTemplate()
		if !ok {
			return "Confirm Scaffold", "The selected template could not be resolved.", errorStyle.Render("Template not found"), "N Back • ESC Quit"
		}
		title := m.appTitle.Value()
		if title == "" {
			title = deriveTitle(m.projectName.Value())
		}
		return "Confirm Scaffold", "Review the scaffold plan before creating files.", m.renderConfirm(tmpl, title), "Y Scaffold • N Restart • ESC Quit"

	case stepScaffolding:
		tmpl, _ := m.selectedTemplate()
		body := fmt.Sprintf("%s Applying %s into %s\n\n%s",
			m.spinner.View(),
			selectedStyle.Render(tmpl.Label),
			codeStyle.Render(filepath.Join(".", m.projectName.Value())),
			subtleStyle.Render("Cloning template, applying manifest rewrites, pruning disabled features, and creating the initial commit."))
		return "Scaffolding", "Your project is being assembled.", body, "Please wait • CTRL+C Quit"

	case stepDone:
		if m.err != nil {
			body := errorStyle.Render("Scaffold failed") + "\n\n" + m.err.Error()
			return "Scaffold Failed", "No successful scaffold result was produced.", body, "Enter Exit • Q Exit"
		}
		return "Project Scaffolded", "The Solar template is ready to use.", m.renderSuccess(), "Enter Exit • Q Exit"
	}

	return "Sun", "Solar project scaffolding.", "", "ESC Quit"
}

func (m model) renderConfirm(tmpl registry.Template, title string) string {
	var b strings.Builder
	b.WriteString(sectionTitle("Project"))
	fmt.Fprintf(&b, "  Name:       %s\n", selectedStyle.Render(m.projectName.Value()))
	if m.projectKind == kindAPI {
		fmt.Fprintf(&b, "  API:        %s\n", m.apiVersion.Value())
		fmt.Fprintf(&b, "  Title:      %s\n", title)
	} else {
		fmt.Fprintf(&b, "  Site title: %s\n", title)
	}
	b.WriteString("\n")
	b.WriteString(sectionTitle("Template"))
	fmt.Fprintf(&b, "  Label:      %s\n", tmpl.Label)
	fmt.Fprintf(&b, "  Ref:        %s\n", tmpl.DefaultRef)
	fmt.Fprintf(&b, "  Repo:       %s\n", tmpl.Repo)
	b.WriteString("\n")
	b.WriteString(sectionTitle("Features"))
	if m.projectKind == kindAPI {
		fmt.Fprintf(&b, "  Google OIDC: %s\n", statusText(m.includeOIDC))
	} else {
		fmt.Fprintf(&b, "  French i18n: %s\n", statusText(m.includeI18n))
		fmt.Fprintf(&b, "  Search:      %s\n", statusText(m.includeSearch))
	}
	b.WriteString("\n")
	b.WriteString(selectedStyle.Render("Press y to create the project."))
	return b.String()
}

func (m model) renderSuccess() string {
	var b strings.Builder
	b.WriteString(successStyle.Render("Project scaffolded successfully."))
	b.WriteString("\n\n")
	b.WriteString(sectionTitle("Result"))
	fmt.Fprintf(&b, "  Location: %s\n", codeStyle.Render(m.result.Destination))
	fmt.Fprintf(&b, "  Template: %s\n", m.result.TemplateID)
	fmt.Fprintf(&b, "  Commit:   %s\n", m.result.InitialCommit)
	if len(m.result.PrunedFeatures) > 0 {
		fmt.Fprintf(&b, "  Pruned:   %s\n", strings.Join(m.result.PrunedFeatures, ", "))
	}
	b.WriteString("\n")
	b.WriteString(sectionTitle("Next commands"))
	b.WriteString(commandBlock(m.nextCommands()))
	return b.String()
}

func (m model) nextCommands() []string {
	commands := []string{fmt.Sprintf("cd %s", m.projectName.Value())}
	if m.projectKind == kindAPI {
		return append(commands, "cp .env.dist .env", "docker compose up")
	}
	return append(commands, "pnpm install", "pnpm start")
}

func (m model) renderFrame(title, subtitle, body, footer string) string {
	width := m.contentWidth()
	header := brandStyle.Render("SUN") + subtleStyle.Render("  Solar Stack Project Forge")
	if current, total := m.stepProgress(); current > 0 && total > 0 {
		header += subtleStyle.Render(fmt.Sprintf("  |  Step %d of %d", current, total))
	} else if current > 0 {
		header += subtleStyle.Render(fmt.Sprintf("  |  Step %d", current))
	}

	content := titleStyle.Render(title)
	if subtitle != "" {
		content += "\n" + subtleStyle.Render(subtitle)
	}
	if body != "" {
		content += "\n\n" + body
	}

	frame := panelStyle.Width(width).Render(content)
	if m.step == stepProjectType || m.step == stepMigrationChoice || m.step == stepConfirm {
		frame = activePanelStyle.Width(width).Render(content)
	}

	if footer != "" {
		footer = "\n" + subtleStyle.Render(footer)
	}
	return docStyle.Render(header + "\n\n" + frame + footer)
}

func (m model) contentWidth() int {
	if m.width <= 0 {
		return 72
	}
	width := m.width - 6
	if width > 86 {
		width = 86
	}
	if width < 32 {
		width = 32
	}
	return width
}

func (m model) inputWidth() int {
	width := m.contentWidth() - 8
	if width > 40 {
		width = 40
	}
	if width < 20 {
		width = 20
	}
	return width
}

func (m model) stepProgress() (int, int) {
	switch m.step {
	case stepProjectType:
		return 1, 0
	case stepMigrationChoice:
		return 2, 7
	case stepProjectName:
		if m.projectKind == kindDocs {
			return 2, 6
		}
		return 3, 7
	case stepAPIVersion:
		return 4, 7
	case stepAppTitle:
		if m.projectKind == kindAPI {
			return 5, 7
		}
		return 3, 6
	case stepOIDC:
		return 6, 7
	case stepI18n:
		return 4, 6
	case stepSearch:
		return 5, 6
	case stepConfirm:
		if m.projectKind == kindAPI {
			return 7, 7
		}
		return 6, 6
	}
	return 0, 0
}

type choiceOption struct {
	Key         string
	Title       string
	Description string
	Meta        string
}

func projectTypeOptions() []choiceOption {
	return []choiceOption{
		{Key: "a", Title: "API backend", Description: "FastAPI service with PostgreSQL, Redis, JWT auth, and migrations.", Meta: "Mars or Mercury"},
		{Key: "d", Title: "Docs site", Description: "Docusaurus documentation site with Tailwind CSS and optional i18n/search.", Meta: "Venus"},
	}
}

func migrationOptions() []choiceOption {
	return []choiceOption{
		{Key: "a", Title: "Alembic", Description: "Recommended Python migration workflow for most API projects.", Meta: "Mars"},
		{Key: "f", Title: "Flyway", Description: "Classic SQL migration workflow for teams that prefer SQL-first changes.", Meta: "Mercury"},
	}
}

func renderChoiceCards(options []choiceOption, cursor int) string {
	var b strings.Builder
	for i, option := range options {
		if i > 0 {
			b.WriteString("\n\n")
		}
		prefix := "  "
		style := optionStyle
		if i == cursor {
			prefix = selectedStyle.Render("> ")
			style = activeOptionStyle
		}
		content := fmt.Sprintf("%s%s %s\n%s%s",
			prefix,
			labelStyle.Render(option.Title),
			keyStyle.Render("("+option.Key+")"),
			"  ",
			subtleStyle.Render(option.Description))
		if option.Meta != "" {
			content += "\n  " + selectedStyle.Render(option.Meta)
		}
		b.WriteString(style.Render(content))
	}
	return b.String()
}

func renderBooleanChoice(question, description string, yes bool) string {
	var b strings.Builder
	b.WriteString(labelStyle.Render(question))
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(description))
	b.WriteString("\n\n")
	b.WriteString(booleanPill("Yes", yes))
	b.WriteString("  ")
	b.WriteString(booleanPill("No", !yes))
	return b.String()
}

func booleanPill(label string, active bool) string {
	style := optionStyle.Padding(0, 2)
	prefix := "  "
	if active {
		style = activeOptionStyle.Padding(0, 2)
		prefix = "> "
	}
	return style.Render(prefix + label)
}

func formField(label, input, hint string) string {
	return labelStyle.Render(label) + "\n" + input + "\n\n" + subtleStyle.Render(hint)
}

func sectionTitle(title string) string {
	return selectedStyle.Render(title) + "\n"
}

func statusText(enabled bool) string {
	if enabled {
		return successStyle.Render("included")
	}
	return subtleStyle.Render("pruned")
}

func commandBlock(commands []string) string {
	var b strings.Builder
	for _, command := range commands {
		b.WriteString(codeStyle.Render(command))
		b.WriteString("\n")
	}
	return b.String()
}

func (m *model) selectedTemplate() (registry.Template, bool) {
	if m.projectKind == kindDocs {
		return registry.ByID("venus")
	}
	return registry.APIByMigration(m.useAlembic)
}

func (m model) updateProjectType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selectDocs := func() {
		m.projectKind = kindDocs
		m.templateID = "venus"
		m.includeI18n = true
		m.i18nChosen = true
		m.includeSearch = true
		m.searchChosen = true
		m.step = stepProjectName
		m.projectName.Focus()
	}
	switch msg.String() {
	case "up", "left", "k", "h":
		m.projectCursor = 0
	case "down", "right", "j", "l":
		m.projectCursor = 1
	case "a", "A":
		m.projectCursor = 0
		m.projectKind = kindAPI
		m.step = stepMigrationChoice
	case "d", "D":
		m.projectCursor = 1
		selectDocs()
	case "enter":
		if m.projectCursor == 0 {
			m.projectKind = kindAPI
			m.step = stepMigrationChoice
			return m, nil
		}
		selectDocs()
	}
	return m, nil
}

func (m model) updateMigrationChoice(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "left", "k", "h":
		m.migrationCursor = 0
	case "down", "right", "j", "l":
		m.migrationCursor = 1
	case "a", "A":
		m.migrationCursor = 0
		m.useAlembic = true
		m.templateID = "mars"
		m.step = stepProjectName
		m.projectName.Focus()
	case "f", "F":
		m.migrationCursor = 1
		m.useAlembic = false
		m.templateID = "mercury"
		m.step = stepProjectName
		m.projectName.Focus()
	case "enter":
		m.useAlembic = m.migrationCursor == 0
		if m.useAlembic {
			m.templateID = "mars"
		} else {
			m.templateID = "mercury"
		}
		m.step = stepProjectName
		m.projectName.Focus()
	}
	return m, nil
}

func (m model) updateProjectName(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.projectName, cmd = m.projectName.Update(msg)
	if msg.String() != "enter" {
		m.validationErr = ""
		return m, cmd
	}
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

func parseBooleanChoice(msg tea.KeyMsg, current bool) (value, chosen bool, advance step) {
	switch msg.String() {
	case "left", "up", "h", "k":
		return true, true, 0
	case "right", "down", "l", "j":
		return false, true, 0
	case "y", "Y":
		return true, true, 0
	case "n", "N":
		return false, true, 0
	case "enter":
		return current, true, -1
	}
	return current, false, 0
}

func (m model) updateOIDC(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	val, chosen, adv := parseBooleanChoice(msg, m.includeOIDC)
	m.includeOIDC = val
	m.oidcChosen = m.oidcChosen || chosen
	if adv == -1 {
		m.step = stepConfirm
	}
	return m, nil
}

func (m model) updateI18n(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	val, chosen, adv := parseBooleanChoice(msg, m.includeI18n)
	m.includeI18n = val
	m.i18nChosen = m.i18nChosen || chosen
	if adv == -1 {
		m.step = stepSearch
	}
	return m, nil
}

func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	val, chosen, adv := parseBooleanChoice(msg, m.includeSearch)
	m.includeSearch = val
	m.searchChosen = m.searchChosen || chosen
	if adv == -1 {
		m.step = stepConfirm
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
