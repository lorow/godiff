package app

import (
	"fmt"
	Shared "godiff/pkg/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m Model) Init() tea.Cmd {
	return m.initScreen(m.screen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debug("Root message received", "Type", fmt.Sprintf("%T", msg), "Message", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
		m.updateScreenSize(m.screen, msg)
		return m, nil

	case Shared.RouteMsg:
		m.screen = msg.To
		var cmds []tea.Cmd
		cmds = append(cmds, m.initScreen(m.screen))
		m.updateScreenSize(m.screen, m.size)
		cmds = append(cmds, m.handleScreenRoute(m.screen, msg.Payload))
		return m, tea.Batch(cmds...)

	case Shared.ExitMsg:
		// placeholder for getting exit confirmation popup later on
		return m, tea.Quit
	}

	cmd := m.updateScreen(m.screen, msg)
	return m, cmd
}

func (m *Model) initScreen(screen Shared.Screen) tea.Cmd {
	var cmd tea.Cmd
	switch screen {
	case Shared.LandingScreen:
		cmd = m.landingPage.Init()

	case Shared.NewProjectWizardScreen:
		cmd = m.newProjectWizard.Init()
	}
	return cmd
}

func (m *Model) updateScreen(screen Shared.Screen, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch screen {
	case Shared.LandingScreen:
		m.landingPage, cmd = m.landingPage.Update(msg)

	case Shared.NewProjectWizardScreen:
		m.newProjectWizard, cmd = m.newProjectWizard.Update(msg)
	}

	return cmd
}

func (m *Model) updateScreenSize(screen Shared.Screen, msg tea.WindowSizeMsg) {
	switch screen {
	case Shared.LandingScreen:
		m.landingPage.SetSize(msg)
	case Shared.NewProjectWizardScreen:
		m.newProjectWizard.SetSize(msg)
	}
}

func (m *Model) handleScreenRoute(screen Shared.Screen, payload any) tea.Cmd {
	// handle payload for specific screens here later
	return nil
}
