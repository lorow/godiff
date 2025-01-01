package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	current_route string
	views         map[string]tea.Model
	commannd_bar  CommandBarModel
}

func newInitialModel() model {
	views := make(map[string]tea.Model)
	views["/"] = NewLandingPage()

	return model{
		current_route: "/",
		views:         views,
		commannd_bar:  NewCommandBar(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for _, view := range m.views {
			view, _ = view.Update(msg)
			m.views[m.current_route] = view
		}
		return m, nil
	}

	command_bar_model, command_bar_message := m.commannd_bar.Update(msg)
	view_model, message := m.views[m.current_route].Update(msg)

	m.commannd_bar = command_bar_model.(CommandBarModel)
	m.views[m.current_route] = view_model

	return m, tea.Batch(command_bar_message, message)
}

func (m model) View() string {
	s := ""

	if current_view, ok := m.views[m.current_route]; ok {
		s += current_view.View()
	}
	s += m.commannd_bar.View()

	return s
}

func main() {
	p := tea.NewProgram(newInitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong while trying to run the program: %v \n", err)
		os.Exit(1)
	}
}
