package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type SetNewSizeMsg struct {
	width, height int
}

type model struct {
	size         tea.WindowSizeMsg
	logs         io.Writer
	currentRoute string
	views        map[string]tea.Model
}

func newInitialModel(logs_file io.Writer) model {
	views := make(map[string]tea.Model)
	views["/"] = NewLandingPage()

	return model{
		logs:         logs_file,
		currentRoute: "/",
		views:        views,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, view := range m.views {
		cmds = append(cmds, view.Init())
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.logs != nil {
		spew.Fdump(m.logs, msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg

		for key, view := range m.views {
			viewModel, _ := view.Update(SetNewSizeMsg{width: msg.Width, height: msg.Height - 1})
			m.views[key] = viewModel
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	viewModel, viewMessage := m.views[m.currentRoute].Update(msg)
	m.views[m.currentRoute] = viewModel

	return m, viewMessage
}

func (m model) View() string {
	appDocument := lipgloss.NewStyle().Width(m.size.Width).Height(m.size.Height)
	doc := strings.Builder{}

	if currentView, ok := m.views[m.currentRoute]; ok {
		doc.WriteString(currentView.View())
	}

	return appDocument.Render(doc.String())
}

// change of plans.
// here's the new idea:
// inspired by the posting.sh thing
// I wanna get rid of the bottom command bar
// instead it's gonna live in a popup - figure out how to make popup with bubble tea
// this popup will have a search input that's gonna go through all of the commands
// and dispatch the proper one once selected.

// I'd also want to have a jump mechanism where by inputing jump button
// we go into jump state, in that state we can cancel or jump focus to something

// landing page view
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   - Projects - 2 -----------------------------------------------------------------------------------------------   |
//|   |                                  				                                                    	   |   |
//|   |  Project name                                 	                                                    	   |   |
//|   |    Short project description                   	                                                    	   |   |
//|   |                                  				                                                    	   |   |
//|   |  Project name                                 	                                                    	   |   |
//|   |    Short project description                   	                                                    	   |   |
//|   |                                  				                                                    	   |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                                                    |
//| up/down select project ^n new project ^o commands enter - load                                                     |
//|--------------------------------------------------------------------------------------------------------------------|

// project page view - single editor
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   |------------------------------------------------------------------------------------------------------------|   |
//|   | GET | http://some-service.dev/                                                                  |          |   |
//|   |------------------------------------------------------------------------------------------------------------|   |
//|                                                                                                                    |
//|   |------------------------------------------------------------------------------------------------------------|   |
//|   |    Short project description                                                                               |   |
//|   |                                                                                                            |   |
//|   |  Project name                                                                                              |   |
//|   |    Short project description                                                                               |   |
//|   |                                                                                                            |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                                                    |
//| <-/up/down/-> change focus ^s save ^o commands ^p jump i edit                                                      |
//|--------------------------------------------------------------------------------------------------------------------|

// project page view - double editor
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   |-----------------------------------------------------|    |-------------------------------------------------|   |
//|   | GET | http://some-service.dev/                |  	|    | GET | http://some-service.dev/            |     |   |
//|   |-----------------------------------------------------|    |-------------------------------------------------|   |
//|                                                                                                                    |
//|   |-----------------------------------------------------|    |-------------------------------------------------|   |
//|   |    Short project description                        |    |    Short project description                    |   |
//|   |                                                     |    |                                                 |   |
//|   |  Project name                                       |    |  Project name                                   |   |
//|   |    Short project description                        |    |    Short project description                    |   |
//|   |                                                     |    |                                                 |   |
//|   -------------------------------------------------------    ---------------------------------------------------   |
//|                                                                                                                    |
//| ^c exit ^s save ^o commands ^p jump i edit enter send                                                              |
//|--------------------------------------------------------------------------------------------------------------------|

// command popup
//|--------------------------------------------------------------------------------------------------------------------|
//|                                                                                                                    |
//|                           |------------------------------------------------------------|                           |
//|                           | Search for command                                         |                           |
//|                           |------------------------------------------------------------|                           |
//|                           |  Some command                                              |                           |
//|                           |  Some command explanation                                  |                           |
//|                           |                                                            |                           |
//|                           |  Some command                                              |                           |
//|                           |  Some command explanation                                  |                           |
//|                           |                                                            |                           |
//|                           |  Some command                                              |                           |
//|                           |  Some command explanation                                  |                           |
//|                           |                                                            |                           |
//|                           |------------------------------------------------------------|                           |
//|                                                                                                                    |
//|                                                                                                                    |
//|--------------------------------------------------------------------------------------------------------------------|

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
	err = InitDatabase()
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
