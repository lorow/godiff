package TitlePanel

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container lipgloss.Style
	Title     lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.Container = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingBottom(1)
	s.Title = lipgloss.NewStyle()

	return s
}
