package views

import (
	"fmt"
	"godiff/components/ItemList"
	"godiff/components/Router"
	"godiff/components/ShortcutsPanel"
	"godiff/components/TitlePanel"
	"godiff/db"
	"godiff/messages"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LandingPageState int

const (
	Initial LandingPageState = iota
	LoadingProjects
	LoadedProjects
	ErrorLoadingProjects
)

type LandingPageModel struct {
	width                   int
	height                  int
	state                   LandingPageState
	itemList                ItemList.Model
	titlePanel              *TitlePanel.Model
	shortcutsPanel          *ShortcutsPanel.Model
	basicShortcuts          []ShortcutsPanel.Shortcut
	onProjectSelectShortcus []ShortcutsPanel.Shortcut
	cursor                  int
	selected                int
}

type SelectedProject int

type RenderableProject struct {
	db.Project
}

func (p RenderableProject) Title() string {
	return fmt.Sprintf("%d - %s", p.Project.ID, p.Project.Name)
}

func (p RenderableProject) Description() string {
	return "Dummy description"
}

func NewLandingPage() LandingPageModel {
	itemRender := ItemList.NewDefaultItemRenderer()
	itemList, _ := ItemList.New("Projects", "No projects loaded", []ItemList.Item{}, itemRender, 1)
	titlePanel := TitlePanel.New(TitlePanel.WithTitle("Welcome to GoDiff - 1.0.0"))

	basicShortcus := []ShortcutsPanel.Shortcut{
		{Key: "^Q", Description: "Quit"},
		{Key: "^O", Description: "Jump focus"},
	}

	onProjectSelectShortcus := []ShortcutsPanel.Shortcut{
		{Key: "^Q", Description: "Quit"},
		{Key: "Enter", Description: "Launch project"},
		{Key: "^O", Description: "Jump focus"},
	}

	shortcutsPanel := ShortcutsPanel.New(
		ShortcutsPanel.WithShortcuts(onProjectSelectShortcus),
	)

	return LandingPageModel{
		itemList:                itemList,
		titlePanel:              titlePanel,
		shortcutsPanel:          shortcutsPanel,
		basicShortcuts:          basicShortcus,
		onProjectSelectShortcus: onProjectSelectShortcus,
		selected:                -1,
	}
}

func (m LandingPageModel) Init() tea.Cmd {
	if m.state == Initial {
		m.state = LoadingProjects
		return LoadProjectsCmd(10, 0)
	}

	return nil
}

func (m LandingPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// we need to first handle the incoming messages
	// only then we can send our own IO commands
	switch msg := msg.(type) {
	case messages.SetNewSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.itemList.SetWidth(msg.Width - 3)
		m.itemList.SetHeight(msg.Height - 10)
		m.shortcutsPanel.SetWidth(msg.Width)

		return m, nil
	case LoadedProjectsMsg:
		m.state = LoadedProjects
		m.itemList.SetItems([]ItemList.Item(msg.projects))
		m.itemList.SetTitle(fmt.Sprintf("Projects - %d", m.itemList.GetItemsCount()))
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.itemList.CursorUp()
			return m, nil
		case "down":
			m.itemList.CursorDown()
			return m, nil
		case "enter":
			m.selected = m.cursor
			currentSelection, _ := m.itemList.GetCurrentSelection()
			return m, Router.RouteTo("editor", SelectedProject(currentSelection.(db.Project).ID))
		}
	}

	return m, nil
}

func (m LandingPageModel) View() string {
	windowContainer := lipgloss.NewStyle().Width(m.width).Height(m.height)
	doc := strings.Builder{}

	title := m.titlePanel.View()
	quitText := lipgloss.NewStyle().PaddingRight(2).Render("Press Q to quit")
	middleSpacer := lipgloss.NewStyle().Width(m.width - lipgloss.Width(title) - lipgloss.Width(quitText)).Render("")
	renderedTitle := lipgloss.NewStyle().PaddingTop(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, title, middleSpacer, quitText))

	doc.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			renderedTitle,
			m.itemList.View(),
			m.shortcutsPanel.View(),
		),
	)

	return windowContainer.Render(doc.String())
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
