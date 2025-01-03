package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CommardBarState int

const (
	CommardBarStateNormal  CommardBarState = iota
	CommandBarStateCommand                 // special state where we're getting a command input
	CommardBarStateEdit
	CommardBarStateVisual
)

var CommandBarStateName = map[CommardBarState]string{
	CommardBarStateNormal: "Normal",
	CommardBarStateEdit:   "Edit",
	CommardBarStateVisual: "Visual",
}

func (s CommardBarState) String() string {
	return CommandBarStateName[s]
}

type CommandBarModel struct {
	// state has dual purpose
	// for one, it'll indicate what's the state of the editor
	// and on which line are we
	// but also it'll indicate if the command bar is selected
	// if we're in CommandBarStateCommand state, we're selected
	// and we're hogging all the input
	// until we are either put out of that state
	// or the command is submitted
	state       CommardBarState
	editor_line [2]int
}

func NewCommandBar() CommandBarModel {
	return CommandBarModel{
		state:       CommardBarStateNormal,
		editor_line: [2]int{0, 0},
	}
}

func (m CommandBarModel) Init() tea.Cmd {
	return nil
}

func (m CommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg == nil {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.state == CommardBarStateNormal {
				return m, tea.Quit
			}
		case ":":
			if m.state == CommardBarStateNormal {
				m.SetState(CommandBarStateCommand)
				return m, nil
			}
		case "v":
			if m.state == CommardBarStateNormal {
				m.SetState(CommardBarStateVisual)
				return m, nil
			}
		case "esc":
			if m.state != CommardBarStateNormal {
				m.SetState(CommardBarStateNormal)
				return m, nil
			}
		}
	}

	return m, nil
}

func (m CommandBarModel) View() string {
	return fmt.Sprintf("commandBar state: %s", m.state)
}

func (m *CommandBarModel) SetState(state CommardBarState) {
	m.state = state
}

func (m CommandBarModel) GetState() CommardBarState {
	return m.state
}
