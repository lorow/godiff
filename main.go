package main

import (
	"fmt"
	"godiff/components/Router"
	"godiff/db"
	"godiff/messages"
	"godiff/views/LandingPage"
	"godiff/views/NewProjectWizard"
	"godiff/views/Shared"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type model struct {
	size   tea.WindowSizeMsg
	router Router.Model
}

func newInitialModel() model {
	view_router := Router.New(
		Router.WithRegisterRoute("/", LandingPage.New()),
		Router.WithRegisterRoute("new_project", NewProjectWizard.New()),
	)

	return model{
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
	log.Debug("Root message received", "Type", fmt.Sprintf("%T", msg), "Message", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
		for key, view := range m.router.GetViews() {
			viewModel, _ := view.Update(messages.SetNewSizeMsg{Width: msg.Width, Height: msg.Height})
			m.router.UpdateRoute(key, viewModel)
		}
		return m, nil
	case Router.RouteMsg:
		return m, m.router.HandleRouteTo(msg)
	case Shared.ExitMsg:
		// placeholder for getting exit confirmation popup later on
		return m, tea.Quit
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
// - I need to handle shortcut inputs - done
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   | Search for a project                                                                                       |   |
//|                                                                                                                    |
//|   - Projects - 2 -----------------------------------------------------------------------------------------------   |
//|   |                                  				                                                    	             |   |
//|   |  Project name                                 	                                                    	     |   |
//|   |    Short project description                   	                                                           |   |
//|   |                                  				                                                                   |   |
//|   |  Project name                                 	                                                   	       |   |
//|   |    Short project description                   	                                                   	       |   |
//|   |                                  				                                                                   |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                                                    |
//| up/down select project ^n new project ^o commands enter - load                                                     |
//|--------------------------------------------------------------------------------------------------------------------|

// new project wizard - will be used to only setup a new project
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   | New Project Title                                                                                          |   |
//|                                                                                                                    |
//|   - Description  -----------------------------------------------------------------------------------------------   |
//|   |                                  				                                                    	             |   |
//|   |                                                                                                            |   |
//|   |                                  				                                                                   |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                         | cancel |  | continue |   |
//| up/down change field esc go back                                                                                   |
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
	f, err := os.OpenFile("debug.log", os.O_TRUNC|os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)

	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	log.SetOutput(f)
	log.SetPrefix("debug")
	defer f.Close()

	// reset the log file after every start, just so we don't accidentally put hundreds of GBs of logs while developing
	f.Truncate(0)
	f.Seek(0, 0)

	log.SetReportCaller(true)
	log.SetReportTimestamp(true)
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)

	log.Info("Starting application...")
	err = db.InitDatabase()
	if err != nil {
		log.Error("Something went wrong while trying to open the database: %v \n", err)
		os.Exit(1)
	}
	log.Debug("Database initialized")

	err = db.MigrateDatabase()
	if err != nil {
		log.Error("Something went wrong while trying to migrate the database: %v \n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(newInitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Error("Something went wrong while trying to run the program: %v \n", err)
		os.Exit(1)
	}
}
