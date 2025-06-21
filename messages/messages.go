package messages

import tea "github.com/charmbracelet/bubbletea"

type SetNewSizeMsg struct {
	Width, Height int
}

type SwitchFocusMsg struct {
	Target string
}

func SwitchFocusCmd(target string) tea.Cmd {
	return func() tea.Msg { return SwitchFocusMsg{Target: target} }
}
