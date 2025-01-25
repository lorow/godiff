package ItemList

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	TitleBar lipgloss.Style
	Title    lipgloss.Style
	Spinner  lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 1, 2) //nolint:mnd

	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#8E8E8E", Dark: "#747373"})
	return s
}
