package BorderWithText

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	BorderStyle lipgloss.Style
	TextStyle   lipgloss.Style
}

// make the apply style function take the style
func DefaultStyles() (s Styles) {
	s.BorderStyle = lipgloss.NewStyle()
	s.TextStyle = lipgloss.NewStyle()

	return s
}

func FocusedStyles() (s Styles) {
	s.BorderStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("62"))
	s.TextStyle = lipgloss.NewStyle()

	return s
}
