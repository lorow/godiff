package FocusChain

import tea "github.com/charmbracelet/bubbletea"

type FocusDirection int

const (
	FocusUp FocusDirection = iota
	FocusDown
)

type SwitchFocusMsg struct {
	Direction FocusDirection
}

func SwitchFocusCmd(direction FocusDirection) tea.Cmd {
	return func() tea.Msg { return SwitchFocusMsg{direction} }
}
