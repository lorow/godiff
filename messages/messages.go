package messages

import tea "github.com/charmbracelet/bubbletea"

type SetNewSizeMsg struct {
	Width, Height int
}

type SwtichFocusMsg struct {
	Target string
}

func SwtichFocusCmd(target string) tea.Cmd {
	return func() tea.Msg { return SwtichFocusMsg{Target: target} }
}
