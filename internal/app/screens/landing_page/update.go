package landingPage

import (
	"fmt"
	"godiff/internal/db"
	Shared "godiff/pkg/common"
	FocusChain "godiff/pkg/components/FocusChain"
	"godiff/pkg/components/ItemList"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	if m.state == Initial {
		m.state = LoadingProjects
		cmds = append(cmds, LoadProjectsCmd(10, 0))
	}

	cmds = append(cmds, m.searchInput.Init())
	return tea.Batch(cmds...)
}

func (m *Model) SetSize(size tea.WindowSizeMsg) {
	m.width = size.Width
	m.height = size.Height
	m.searchInput.SetWidth(size.Width - 3)
	m.itemList.SetWidth(size.Width - 3)
	m.itemList.SetHeight(size.Height - 10)
	m.shortcutsPanel.SetWidth(size.Width)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// we need to first handle the incoming messages
	// only then we can send our own IO commands
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case LoadedProjectsMsg:
		m.state = LoadedProjects
		m.itemList.SetItems([]ItemList.Item(msg.projects))
		if len(msg.projects) != 0 {
			m.itemList.SetTitle(fmt.Sprintf("%d Projects loaded - choose a project to work or create a new one", m.itemList.GetItemsCount()))
		} else {
			m.itemList.SetTitle("No projects found, create a new one!")
		}
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

	currentlySelected := m.focusChain.GetCurrentlySelected()
	if currentlySelected == m.searchInput {
		m.shortcutsPanel.SetShortcuts(m.basicShortcuts)
	}

	if currentlySelected == m.itemList {
		m.shortcutsPanel.SetShortcuts(m.onProjectSelectShortcuts)
	}

	cmds = append(cmds, m.shortcutsPanel.CheckIfShortcutHit(msg))
	cmds = append(cmds, currentlySelected.Update(msg))
	return m, tea.Batch(cmds...)
}

func onItemListSelect(item ItemList.Item) tea.Cmd {
	return Shared.RouteTo(Shared.ProjectEditorScreen, SelectedProject(item.(RenderableProject).ID))
}

func onSearchSubmit(input string) tea.Cmd {
	return Shared.FilterProjectsCmd(input)
}

type LoadedProjectsMsg struct {
	projects []ItemList.Item
}

func LoadProjectsCmd(limit, offset int) tea.Cmd {
	return func() tea.Msg {
		var projects []ItemList.Item

		for _, project := range db.GetProjects(limit, offset) {
			projects = append(projects, RenderableProject{project})
		}

		return LoadedProjectsMsg{projects}
	}
}
