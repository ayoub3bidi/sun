package tui

import (
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
	stepProjectName step = iota
	stepAPIVersion
	stepAppTitle
	stepOIDC
	stepConfirm
	stepScaffolding
	stepDone
)

type scaffoldResultMsg struct {
	result sc.ScaffoldResult
	err    error
}

type model struct {
	step        step
	projectName textinput.Model
	apiVersion  textinput.Model
	appTitle    textinput.Model
	includeOIDC bool
	oidcChosen  bool
	spinner     spinner.Model
	err         error
	result      *sc.ScaffoldResult
	width       int
	height      int
}

func NewScaffoldModel() model {
	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Dot

	pn := textinput.New()
	pn.Placeholder = "my-api"
	pn.Focus()
	pn.CharLimit = 64
	pn.Width = 40

	av := textinput.New()
	av.Placeholder = "v1"
	av.CharLimit = 16
	av.Width = 40

	at := textinput.New()
	at.Placeholder = "My API"
	at.CharLimit = 64
	at.Width = 40

	return model{
		step:        stepProjectName,
		projectName: pn,
		apiVersion:  av,
		appTitle:    at,
		spinner:     s,
	}
}

func Run(m model) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
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
		case stepProjectName:
			return m.updateProjectName(msg)
		case stepAPIVersion:
			return m.updateAPIVersion(msg)
		case stepAppTitle:
			return m.updateAppTitle(msg)
		case stepOIDC:
			return m.updateOIDC(msg)
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
			m.err = msg.err
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

	switch m.step {
	case stepProjectName:
		b.WriteString(titleStyle.Render("New Solar Project"))
		b.WriteString("\n\n")
		b.WriteString(labelStyle.Render("Project name:"))
		b.WriteString("\n")
		b.WriteString(m.projectName.View())
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Kebab-case name, e.g. my-api"))

	case stepAPIVersion:
		b.WriteString(titleStyle.Render("API Version Prefix"))
		b.WriteString("\n\n")
		b.WriteString(labelStyle.Render("API version prefix:"))
		b.WriteString("\n")
		b.WriteString(m.apiVersion.View())
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("e.g. v1, v2, v2026-01"))

	case stepAppTitle:
		b.WriteString(titleStyle.Render("Application Title"))
		b.WriteString("\n\n")
		b.WriteString(labelStyle.Render("App title:"))
		b.WriteString("\n")
		b.WriteString(m.appTitle.View())
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Used in page titles and documentation"))

	case stepOIDC:
		b.WriteString(titleStyle.Render("OIDC Authentication"))
		b.WriteString("\n\n")
		b.WriteString("Include Google OIDC?")
		b.WriteString("\n\n")
		var chosen string
		if m.includeOIDC {
			chosen = "[✓] Yes"
		} else if m.oidcChosen {
			chosen = "[✓] No"
		} else {
			chosen = "[ ]"
		}
		b.WriteString(chosen)
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Press 'y' for Yes, 'n' for No, then Enter to confirm"))

	case stepConfirm:
		tmpl, _ := registry.ByID("mars")
		title := m.appTitle.Value()
		if title == "" {
			title = deriveTitle(m.projectName.Value())
		}
		b.WriteString(titleStyle.Render("Confirm Scaffold"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  Project name:  %s\n", m.projectName.Value()))
		b.WriteString(fmt.Sprintf("  API version:   %s\n", m.apiVersion.Value()))
		b.WriteString(fmt.Sprintf("  App title:     %s\n", title))
		b.WriteString(fmt.Sprintf("  OIDC:          %v\n", m.includeOIDC))
		b.WriteString(fmt.Sprintf("  Template:      %s @ %s\n", tmpl.ID, tmpl.DefaultRef))
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
			b.WriteString("\n")
			b.WriteString("\n")
			b.WriteString(subtleStyle.Render("Next steps:"))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf("  cd %s\n", m.projectName.Value()))
			b.WriteString("  cp .env.dist .env")
			b.WriteString("\n\n")
			b.WriteString(subtleStyle.Render("Press Enter or 'q' to exit"))
		}
	}

	return docStyle.Render(b.String())
}

func (m model) updateProjectName(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.projectName, cmd = m.projectName.Update(msg)
	if msg.String() == "enter" && m.projectName.Value() != "" {
		m.projectName.Blur()
		m.step = stepAPIVersion
		m.apiVersion.Focus()
		m.apiVersion.SetValue("v1")
	}
	return m, cmd
}

func (m model) updateAPIVersion(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.apiVersion, cmd = m.apiVersion.Update(msg)
	if msg.String() == "enter" && m.apiVersion.Value() != "" {
		m.apiVersion.Blur()
		m.step = stepAppTitle
		m.appTitle.Focus()
		title := deriveTitle(m.projectName.Value())
		m.appTitle.SetValue(title)
	}
	return m, cmd
}

func (m model) updateAppTitle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.appTitle, cmd = m.appTitle.Update(msg)
	if msg.String() == "enter" {
		m.appTitle.Blur()
		m.step = stepOIDC
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

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.step = stepScaffolding
		return m, scaffoldCmd(m)
	case "n", "N":
		m.step = stepProjectName
		m.oidcChosen = false
		m.projectName.Focus()
		return m, nil
	}
	return m, nil
}

func scaffoldCmd(m model) tea.Cmd {
	return func() tea.Msg {
		tmpl, _ := registry.ByID("mars")
		title := m.appTitle.Value()
		if title == "" {
			title = deriveTitle(m.projectName.Value())
		}
		dest := filepath.Join(".", m.projectName.Value())
		config := sc.ScaffoldConfig{
			Repo:        tmpl.Repo,
			Ref:         tmpl.DefaultRef,
			Destination: dest,
			Variables: map[string]string{
				"project_name": m.projectName.Value(),
				"api_version":  m.apiVersion.Value(),
				"app_title":    title,
			},
			InitialCommit: "Initial scaffold",
		}
		if !m.includeOIDC {
			config.DisabledFeatures = []string{"oidc"}
		}
		result, err := sc.Scaffold(config)
		return scaffoldResultMsg{result, err}
	}
}

func deriveTitle(s string) string {
	if s == "" {
		return ""
	}
	words := strings.FieldsFunc(s, func(r rune) bool { return r == '-' || r == '_' })
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
