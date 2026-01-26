package Shared

import tea "github.com/charmbracelet/bubbletea"

type FilterProjectsMsg struct {
	query string
}

func FilterProjectsCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return FilterProjectsMsg{query}
	}
}

type JumpFocusEnableMsg struct{}

func JumpFocusCmd() tea.Cmd {
	return func() tea.Msg {
		return JumpFocusEnableMsg{}
	}
}

type ExitMsg struct{}

func ExitCmd() tea.Cmd {
	return func() tea.Msg {
		return ExitMsg{}
	}
}

func RouteTo(screen Screen, payload any) tea.Cmd {
	return func() tea.Msg {
		return RouteMsg{To: screen, Payload: payload}
	}
}
