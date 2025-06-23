package views

import tea "github.com/charmbracelet/bubbletea"

type FilterProjectsMsg struct {
	query string
}

func FilterProjectsCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return FilterProjectsMsg{query}
	}
}
