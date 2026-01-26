package newProjectWizard

import (
	"godiff/pkg/components/FocusChain"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetSize(size tea.WindowSizeMsg) {
	m.width = size.Width
	m.height = size.Height
	m.projectName.SetWidth(size.Width - 3)
	m.projectDescription.SetWidth(size.Width - 3)
	m.shortcutsPanel.SetWidth(size.Width)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case FocusChain.SwitchFocusMsg:
		if msg.Direction == FocusChain.FocusUp {
			cmd, _ := m.focusChain.Previous()
			cmds = append(cmds, cmd)
		}

		if msg.Direction == FocusChain.FocusDown {
			cmd, _ := m.focusChain.Next()
			cmds = append(cmds, cmd)
		}
	}

	cmds = append(cmds, m.shortcutsPanel.CheckIfShortcutHit(msg))
	currentlySelected := m.focusChain.GetCurrentlySelected()
	cmds = append(cmds, currentlySelected.Update(msg))

	return m, tea.Batch(cmds...)
}
