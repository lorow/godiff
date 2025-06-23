package TextInput

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.Container = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).MarginLeft(1).MarginRight(1).PaddingLeft(1)
	return s
}

func focusedStyles() (s Styles) {
	s.Container = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).MarginLeft(1).MarginRight(1).PaddingLeft(1)
	return s
}
