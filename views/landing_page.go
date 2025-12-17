package views

import (
	"fmt"
	"strings"

	"godiff/components/FocusChain"
	"godiff/components/ItemList"
	"godiff/components/Router"
	"godiff/components/ShortcutsPanel"
	"godiff/components/TextInput"
	"godiff/components/TitlePanel"
	"godiff/db"
	"godiff/messages"

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
	width                    int
	height                   int
	state                    LandingPageState
	searchInput              *TextInput.Model
	itemList                 *ItemList.Model
	titlePanel               *TitlePanel.Model
	shortcutsPanel           *ShortcutsPanel.Model
	basicShortcuts           []ShortcutsPanel.Shortcut
	onProjectSelectShortcuts []ShortcutsPanel.Shortcut
	focusChain               *FocusChain.Model
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
	searchInput := TextInput.New(TextInput.WIthOnSubmit(onSearchSubmit))
	searchInput.Focus()

	itemList := ItemList.New(
		ItemList.WithTittle("Projects"),
		ItemList.WithNoItemsText("No projects loaded"),
		ItemList.WithOnSelection(onItemListSelect),
	)

	titlePanel := TitlePanel.New(TitlePanel.WithTitle("Welcome to GoDiff - 1.0.0"))

	basicShortcuts := []ShortcutsPanel.Shortcut{
		{Key: "↑/↓", Description: "Change focus"},
		{Key: "^Q", Description: "Quit"},
		{Key: "^N", Description: "Create new project"},
		{Key: "^O", Description: "Jump focus"},
	}

	onProjectSelectShortcuts := []ShortcutsPanel.Shortcut{
		{Key: "↑/↓", Description: "Change focus"},
		{Key: "^Q", Description: "Quit"},
		{Key: "^N", Description: "Create new project"},
		{Key: "Enter", Description: "Launch project"},
		{Key: "^O", Description: "Jump focus"},
	}

	shortcutsPanel := ShortcutsPanel.New(
		ShortcutsPanel.WithShortcuts(onProjectSelectShortcuts),
	)

	focusChain := FocusChain.New(FocusChain.WithItem(searchInput), FocusChain.WithItem(itemList))

	return LandingPageModel{
		searchInput:              searchInput,
		itemList:                 itemList,
		titlePanel:               titlePanel,
		shortcutsPanel:           shortcutsPanel,
		basicShortcuts:           basicShortcuts,
		onProjectSelectShortcuts: onProjectSelectShortcuts,
		focusChain:               focusChain,
	}
}

func (m LandingPageModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	if m.state == Initial {
		m.state = LoadingProjects
		cmds = append(cmds, LoadProjectsCmd(10, 0))
	}

	cmds = append(cmds, m.searchInput.Init())
	return tea.Batch(cmds...)
}

func (m LandingPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// we need to first handle the incoming messages
	// only then we can send our own IO commands
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case messages.SetNewSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.searchInput.SetWidth(msg.Width - 3)
		m.itemList.SetWidth(msg.Width - 3)
		m.itemList.SetHeight(msg.Height - 10)
		m.shortcutsPanel.SetWidth(msg.Width)

		return m, nil
	case LoadedProjectsMsg:
		m.state = LoadedProjects
		m.itemList.SetItems([]ItemList.Item(msg.projects))
		m.itemList.SetTitle(fmt.Sprintf("%d Projects loaded - choose a project to work or create a new one", m.itemList.GetItemsCount()))

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

	cmds = append(cmds, currentlySelected.Update(msg))
	return m, tea.Batch(cmds...)
}

func onItemListSelect(item ItemList.Item) tea.Cmd {
	return Router.RouteTo("editor", SelectedProject(item.(RenderableProject).ID))
}

func onSearchSubmit(input string) tea.Cmd {
	return FilterProjectsCmd(input)
}

func (m LandingPageModel) View() string {
	windowContainer := lipgloss.NewStyle().Width(m.width).Height(m.height)
	doc := strings.Builder{}

	title := m.titlePanel.View()
	quitText := lipgloss.NewStyle().PaddingRight(2).Render("Press Q to quit")
	middleSpacer := lipgloss.NewStyle().Width(m.width - lipgloss.Width(title) - lipgloss.Width(quitText)).Render("")
	renderedTitle := lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, title, middleSpacer, quitText))

	doc.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			renderedTitle,
			m.searchInput.View(),
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
