package main

import tea "github.com/charmbracelet/bubbletea"

type CommandBarModel struct {
}

func NewCommandBar() CommandBarModel {
	return CommandBarModel{}
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
			return m, tea.Quit
		}
	}

	return m, func() tea.Msg { return msg }
}

func (m CommandBarModel) View() string {
	return ""
}
