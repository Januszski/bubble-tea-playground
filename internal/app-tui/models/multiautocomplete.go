package models

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type MultiAutocomplete struct {
	title    string
	options  []string
	selected []string

	input  textinput.Model
	cursor int
}

func NewMultiAutocomplete(title string, options []string) MultiAutocomplete {
	input := textinput.New()
	input.Placeholder = "Start typing to search..."
	input.Prompt = ""

	return MultiAutocomplete{
		title:   title,
		options: options,
		input:   input,
	}
}

func (m *MultiAutocomplete) Focus() tea.Cmd {
	return m.input.Focus()
}

func (m *MultiAutocomplete) Blur() {
	m.input.Blur()
}

// SelectedValues returns all selected options.
func (m MultiAutocomplete) SelectedValues() []string {
	return m.selected
}

// HasQuery tells the parent whether the user is currently searching.
func (m MultiAutocomplete) HasQuery() bool {
	return m.input.Value() != ""
}

// ClearQuery clears the current search.
func (m *MultiAutocomplete) ClearQuery() {
	m.input.SetValue("")
	m.cursor = 0
}

// filteredOptions returns options matching the current query.
func (m MultiAutocomplete) filteredOptions() []string {
	query := strings.ToLower(m.input.Value())

	if query == "" {
		return nil
	}

	var matches []string

	for _, option := range m.options {
		if strings.Contains(strings.ToLower(option), query) {
			// Don't suggest something already selected.
			if !contains(m.selected, option) {
				matches = append(matches, option)
			}
		}
	}

	return matches
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}

	return false
}

func (m MultiAutocomplete) Focused() bool {
	return m.input.Focused()
}

// Update handles keyboard interaction.
func (m MultiAutocomplete) Update(msg tea.Msg) (MultiAutocomplete, tea.Cmd) {
	var cmd tea.Cmd

	// Not currently typing: wait for Enter to activate.
	if !m.input.Focused() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "enter" {
				cmd = m.input.Focus()
				return m, cmd
			}
		}

		return m, nil
	}

	// Currently typing/searching.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down":
			matches := m.filteredOptions()

			if m.cursor < len(matches)-1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			matches := m.filteredOptions()

			if len(matches) > 0 && m.cursor < len(matches) {
				selected := matches[m.cursor]
				m.selected = append(m.selected, selected)

				m.input.SetValue("")
				m.cursor = 0
			}

			return m, nil

		case "esc":
			m.input.Blur()
			m.input.SetValue("")
			m.cursor = 0
			return m, nil
		}
	}

	// Delegate ordinary typing to textinput.
	oldValue := m.input.Value()

	m.input, cmd = m.input.Update(msg)

	if oldValue != m.input.Value() {
		m.cursor = 0
	}

	return m, cmd
}

// View renders the search input and suggestions.
func (m MultiAutocomplete) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	s.WriteString(titleStyle.Render(m.title))
	s.WriteString("\n\n")

	// Display selected options.
	if len(m.selected) > 0 {
		s.WriteString("Selected:\n")

		for _, option := range m.selected {
			s.WriteString(selectedStyle.Render("  ✓ "+option) + "\n")
		}

		s.WriteString("\n")
	}

	s.WriteString("Add another:\n")
	s.WriteString(m.input.View())
	s.WriteString("\n")

	matches := m.filteredOptions()

	if len(matches) > 0 {
		s.WriteString("\n")

		for i, option := range matches {
			line := fmt.Sprintf("  %s", option)

			if i == m.cursor {
				line = cursorStyle.Render("> " + option)
			}

			s.WriteString(line + "\n")
		}
	} else if m.input.Value() != "" {
		s.WriteString(mutedStyle.Render("  No matching options") + "\n")
	}

	s.WriteString("\n")
	if !m.input.Focused() {
		s.WriteString(
			mutedStyle.Render("Press Enter to start typing"),
		)
	} else {
		s.WriteString(
			mutedStyle.Render("↑/↓ Navigate  Enter Select  Esc Stop typing"),
		)
	}

	return s.String()

}
