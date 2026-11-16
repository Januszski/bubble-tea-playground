package apptui

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Team string
type Group string

type Config struct {
	teams  []Team
	groups []Group
}

type App struct {
	width  int
	height int

	help     help.Model
	loaded   bool
	quitting bool

	configured bool
	config     Config

	wizard   Wizard
	findings []Finding
}

func NewApp() *App {
	help := help.New()
	help.ShowAll = true
	return &App{help: help, configured: false, wizard: NewWizard()}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle application-wide messages first.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	}

	// Delegate to the active screen.
	if !a.configured {
		updated, cmd := a.wizard.Update(msg)
		a.wizard = updated
		return a, cmd
	}

	// Handle normal application messages here.

	return a, nil
}

func (a App) View() tea.View {
	var view tea.View

	view.AltScreen = true

	content := a.wizard.View()

	styled := lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	view.SetContent(styled)

	return view
}
