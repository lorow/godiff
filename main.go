package main

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	logs          io.Writer
	current_route string
	views         map[string]tea.Model
	commannd_bar  CommandBarModel
}

func newInitialModel(logs_file io.Writer) model {
	views := make(map[string]tea.Model)
	views["/"] = NewLandingPage()

	return model{
		logs:          logs_file,
		current_route: "/",
		views:         views,
		commannd_bar:  NewCommandBar(),
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
	var dump *os.File

	if _, ok := os.LookupEnv("GODIFF_DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	p := tea.NewProgram(newInitialModel(dump))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong while trying to run the program: %v \n", err)
		os.Exit(1)
	}
}
