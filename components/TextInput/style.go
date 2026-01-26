package TextInput

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.Container = lipgloss.NewStyle().PaddingLeft(1)
	return s
}

func FocusedStyles() (s Styles) {
	s.Container = lipgloss.NewStyle().BorderForeground(lipgloss.Color("62")).PaddingLeft(1)
	return s
}
