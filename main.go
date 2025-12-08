package main

import (
	"fmt"
	"godiff/components/Router"
	"godiff/db"
	"godiff/messages"
	"godiff/views"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	size   tea.WindowSizeMsg
	logs   io.Writer
	router Router.Model
}

func newInitialModel(logs_file io.Writer) model {
	view_router := Router.New(
		Router.WithStartingPage(views.NewLandingPage()),
	)

	return model{
		logs:   logs_file,
		router: *view_router,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, view := range m.router.GetViews() {
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
		for key, view := range m.router.GetViews() {
			viewModel, _ := view.Update(messages.SetNewSizeMsg{Width: msg.Width, Height: msg.Height})
			m.router.UpdateRoute(key, viewModel)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	currentView := m.router.GetCurrentVIew()
	viewModel, viewMessage := currentView.Update(msg)
	m.router.UpdateRoute(m.router.GetCUrrentRoute(), viewModel)

	return m, viewMessage
}

func (m model) View() string {
	doc := strings.Builder{}
	currentView := m.router.GetCurrentVIew()
	doc.WriteString(currentView.View())

	return doc.String()
}

// change of plans.
// here's the new idea:
// inspired by the posting.sh thing
// I wanna get rid of the bottom command bar
// instead it's gonna live in a popup - figure out how to make popup with bubble tea
// this popup will have a search input that's gonna go through all of the commands
// and dispatch the proper one once selected.

// I'd also want to have a jump mechanism where by inputting jump button
// we go into jump state, in that state we can cancel or jump focus to something

// landing page view
// progress on that part:
// - UI mostly done - need to rethink styles
// - I need to handle the filtering command
// - I need to handle shortcut inputs
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   | Search for a project                                                                                       |   |
//|                                                                                                                    |
//|   - Projects - 2 -----------------------------------------------------------------------------------------------   |
//|   |                                  				                                                    	     |   |
//|   |  Project name                                 	                                                    	     |   |
//|   |    Short project description                   	                                                    	     |   |
//|   |                                  				                                                             |   |
//|   |  Project name                                 	                                                    	     |   |
//|   |    Short project description                   	                                                    	     |   |
//|   |                                  				                                                    	             |   |
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
//|   | GET | http://some-service.dev/                 |  	|    | GET | http://some-service.dev/            |     |   |
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
	err = db.InitDatabase()
	if err != nil {
		fmt.Printf("Something went wrong while trying to open the database: %v \n", err)
		os.Exit(1)
	}

	err = db.MigrateDatabase()
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
