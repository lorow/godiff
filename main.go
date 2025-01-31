package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"godiff/components/ItemList"
	"io"
	"os"
	"strings"

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
	testList     ItemList.Model
	views        map[string]tea.Model
}

type testItem struct{}

func (i testItem) Title() string {
	return "test Title"
}

func (i testItem) Description() string {
	return "Some testing description just to have something to work with"
}

func newInitialModel(logs_file io.Writer) model {
	views := make(map[string]tea.Model)
	views["/"] = NewLandingPage()
	listItems := []ItemList.Item{
		testItem{},
		testItem{},
		testItem{},
	}

	return model{
		logs:         logs_file,
		currentRoute: "/",
		views:        views,
		testList:     ItemList.New("test", "no items loaded", listItems, 1),
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
		m.testList.SetWidth(msg.Width - 3)
		m.testList.SetHeight(msg.Height - 10)
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

	doc.WriteString(m.testList.View())

	//if currentView, ok := m.views[m.currentRoute]; ok {
	//	doc.WriteString(currentView.View())
	//}

	return appDocument.Render(doc.String())
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
