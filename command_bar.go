package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type CommandBarState int

const (
	CommandBarStateNormal  CommandBarState = iota
	CommandBarStateCommand                 // special state where we're getting a command input
	CommandBarStateEdit
	CommandBarStateVisual
)

var CommandBarStateName = map[CommandBarState]string{
	CommandBarStateNormal: "Normal",
	CommandBarStateEdit:   "Edit",
	CommandBarStateVisual: "Visual",
}

func (s CommandBarState) String() string {
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
	state      CommandBarState
	width      int
	height     int
	editorLine [2]int
}

func NewCommandBar() CommandBarModel {
	return CommandBarModel{
		state:      CommandBarStateNormal,
		editorLine: [2]int{0, 0},
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
	case SetNewSizeMsg:
		m.width = msg.width
		m.height = msg.height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.state == CommandBarStateNormal {
				return m, tea.Quit
			}
		case ":":
			if m.state == CommandBarStateNormal {
				m.SetState(CommandBarStateCommand)
				return m, nil
			}
		case "v":
			if m.state == CommandBarStateNormal {
				m.SetState(CommandBarStateVisual)
				return m, nil
			}
		case "esc":
			if m.state != CommandBarStateNormal {
				m.SetState(CommandBarStateNormal)
				return m, nil
			}
		}
	}

	return m, nil
}

func (m CommandBarModel) View() string {
	statusBarStyle := lipgloss.NewStyle().Height(m.height).Width(m.width).Background(lipgloss.Color("240"))

	return statusBarStyle.Render(fmt.Sprintf("commandBar state: %s, width: %d", m.state, m.width))
}

func (m *CommandBarModel) SetState(state CommandBarState) {
	m.state = state
}

func (m CommandBarModel) GetState() CommandBarState {
	return m.state
}
