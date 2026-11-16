package models

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type MultiSelect struct {
	title    string
	options  []string
	selected map[string]bool
	cursor   int
}

func NewMultiSelect(title string, options []string) MultiSelect {
	return MultiSelect{
		title:    title,
		options:  options,
		selected: make(map[string]bool),
	}
}

func (m MultiSelect) SelectedValues() []string {
	var selected []string

	for _, option := range m.options {
		if m.selected[option] {
			selected = append(selected, option)
		}
	}

	return selected
}

func (m MultiSelect) Update(msg tea.Msg) (MultiSelect, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "space":
			if len(m.options) > 0 {
				option := m.options[m.cursor]
				m.selected[option] = !m.selected[option]
			}
		}
	}

	return m, nil
}

func (m MultiSelect) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	s.WriteString(titleStyle.Render(m.title))
	s.WriteString("\n\n")

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checkbox := "[ ]"
		if m.selected[option] {
			checkbox = "[x]"
		}

		line := fmt.Sprintf("%-2s %-4s %s", cursor, checkbox, option)

		if m.cursor == i {
			line = cursorStyle.Render(line)
		} else if m.selected[option] {
			line = selectedStyle.Render(line)
		}

		s.WriteString(line + "\n")
	}

	s.WriteString("\n↑/↓ Navigate  Space Select")

	return s.String()
}
