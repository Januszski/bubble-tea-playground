package apptui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/Januszski/bubble-tea-playground/internal/app-tui/models"
)

type WizardStep int

const (
	WelcomeStep WizardStep = iota
	TeamsStep
	GroupsStep
	ReviewStep
)

type Wizard struct {
	step WizardStep

	githubTeams models.MultiSelect
	appGroups   models.MultiSelect
}

func (w *Wizard) Next() {
	if w.step < ReviewStep {
		w.step++
	}
}

func (w *Wizard) Previous() {
	if w.step > WelcomeStep {
		w.step--
	}
}

func NewWizard() Wizard {
	return Wizard{
		step: WelcomeStep,

		githubTeams: models.NewMultiSelect("Select your GitHub teams", []string{
			"Backend",
			"Frontend",
			"Infrastructure",
			"Security",
			"DevOps",
		}),

		appGroups: models.NewMultiSelect("Select your application groups", []string{
			"Production",
			"Staging",
			"Development",
		}),
	}
}

func (w Wizard) Init() tea.Cmd {
	return nil
}

func (w Wizard) Update(msg tea.Msg) (Wizard, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			w.Previous()
			return w, nil

		case "right", "l":
			w.Next()
			return w, nil

		case "enter":
			if w.step == TeamsStep && len(w.githubTeams.SelectedValues()) == 0 {
				return w, nil
			}

			if w.step == GroupsStep && len(w.appGroups.SelectedValues()) == 0 {
				return w, nil
			}

			w.Next()
			return w, nil

		case "q", "ctrl+c":
			return w, tea.Quit
		}
	}

	var cmd tea.Cmd

	switch w.step {
	case TeamsStep:
		w.githubTeams, cmd = w.githubTeams.Update(msg)

	case GroupsStep:
		w.appGroups, cmd = w.appGroups.Update(msg)
	}

	return w, cmd
}

func (w Wizard) View() string {
	switch w.step {
	case WelcomeStep:
		return "Welcome!\n\nPress right to continue."

	case TeamsStep:
		return w.githubTeams.View() + "\n\nPress enter to continue."

	case GroupsStep:
		return w.appGroups.View() + "\n\nPress enter to continue."

	case ReviewStep:
		return fmt.Sprintf(
			"Review your selections\n\nTeams: %v\nGroups: %v\n\nPress enter to finish.",
			w.githubTeams.SelectedValues(),
			w.appGroups.SelectedValues(),
		)
	}

	return ""
}
