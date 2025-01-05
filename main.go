package main

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type SetNewSizeMsg struct {
	width, height int
}

type model struct {
	logs         io.Writer
	currentRoute string
	views        map[string]tea.Model
	commandBar   CommandBarModel
}

func newInitialModel(logs_file io.Writer) model {
	views := make(map[string]tea.Model)
	views["/"] = NewLandingPage()

	return model{
		logs:         logs_file,
		currentRoute: "/",
		views:        views,
		commandBar:   NewCommandBar(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.logs != nil {
		spew.Fdump(m.logs, msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// technically we should be sending two different messages here
		// one for the views and one for the commandbar
		// but that'd mean we'd have to update **every** view
		// each time we get a message.
		// so this should be fine
		for key, view := range m.views {
			viewModel, _ := view.Update(SetNewSizeMsg{width: msg.Width, height: msg.Height - 1})
			m.views[key] = viewModel
		}
		commandBarModel, _ := m.commandBar.Update(SetNewSizeMsg{width: msg.Width, height: 1})
		m.commandBar = commandBarModel.(CommandBarModel)
		return m, nil
	}

	commandBarModel, commandBarMessage := m.commandBar.Update(msg)
	m.commandBar = commandBarModel.(CommandBarModel)
	// if the command bar returned anything
	// it means it reacted to the current message
	// so we need to return it and react to it in the next update
	if commandBarMessage != nil {
		return m, commandBarMessage
	}

	// if we're in command mode, we can't have views receive any input
	if m.commandBar.GetState() == CommandBarStateCommand {
		return m, nil
	}

	viewModel, viewMessage := m.views[m.currentRoute].Update(msg)
	m.views[m.currentRoute] = viewModel

	return m, viewMessage
}

func (m model) View() string {
	s := ""

	if currentView, ok := m.views[m.currentRoute]; ok {
		s += currentView.View()
	}
	s += m.commandBar.View()

	return s
}

func main() {
	var dump *os.File

	if _, ok := os.LookupEnv("GODIFF_DEBUG"); !ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	var err error
	err = InitiDatabase()
	if err != nil {
		fmt.Printf("Something went wrong while trying to open the database: %v \n", err)
		os.Exit(1)
	}

	err = MigrateDatabase()
	if err != nil {
		fmt.Printf("Something went wrong while trying to migrate the database: %v \n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(newInitialModel(dump), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong while trying to run the program: %v \n", err)
		os.Exit(1)
	}
}
