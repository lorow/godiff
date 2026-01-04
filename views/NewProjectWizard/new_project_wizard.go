package NewProjectWizard

import (
	"godiff/components/FocusChain"
	"godiff/components/Router"
	"godiff/components/ShortcutsPanel"
	"godiff/components/TextInput"
	"godiff/components/TitlePanel"
	"godiff/messages"
	"godiff/views/Shared"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width              int
	height             int
	titlePanel         *TitlePanel.Model
	projectName        *TextInput.Model // todo we need borders with names
	projectDescription *TextInput.Model // make this a proper text area later
	shortcutsPanel     *ShortcutsPanel.Model
	focusChain         *FocusChain.Model
}

func New() Model {
	projectName := TextInput.New()
	projectDescription := TextInput.New()

	titlePanel := TitlePanel.New(TitlePanel.WithTitle("Create a new project!"))

	shortcuts := []ShortcutsPanel.Shortcut{
		ShortcutsPanel.NewShortcut("↑/↓", Shared.NoopShortcut, "Change focus", nil), // display only
		ShortcutsPanel.NewShortcut("esc", Shared.EscapeShortcut, "Go back", Router.RouteTo("/", nil)),
		ShortcutsPanel.NewShortcut("^Q", Shared.ExitShortcut, "Quit", Shared.ExitCmd()),
		ShortcutsPanel.NewShortcut("^O", Shared.JumpFocusShortcut, "Jump focus", Shared.JumpFocusCmd()),
	}

	shortcutsPanel := ShortcutsPanel.New(ShortcutsPanel.WithShortcuts(shortcuts))

	focusChain := FocusChain.New(FocusChain.WithItem(projectName), FocusChain.WithItem(projectDescription))

	return Model{
		titlePanel:         titlePanel,
		focusChain:         focusChain,
		shortcutsPanel:     shortcutsPanel,
		projectName:        projectName,
		projectDescription: projectDescription,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case messages.SetNewSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.projectName.SetWidth(msg.Width - 3)
		m.projectDescription.SetWidth(msg.Width - 3)
		m.shortcutsPanel.SetWidth(msg.Width)
		return m, nil
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

	cmds = append(cmds, m.shortcutsPanel.CheckIfShortcutHit(msg))
	currentlySelected := m.focusChain.GetCurrentlySelected()
	cmds = append(cmds, currentlySelected.Update(msg))

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	windowContainer := lipgloss.NewStyle().Width(m.width).Height(m.height)
	doc := strings.Builder{}

	title := m.titlePanel.View()
	middleSpacer := lipgloss.NewStyle().Width(m.width - lipgloss.Width(title)).Render("")
	renderedTitle := lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, title, middleSpacer))

	projectNameView := m.projectName.View()
	projectDescriptionView := m.projectDescription.View()
	shortcutsView := m.shortcutsPanel.View()

	spaceFiller := lipgloss.NewStyle().Width(m.width).Height(m.height - lipgloss.Height(renderedTitle) - lipgloss.Height(projectNameView) - lipgloss.Height(projectDescriptionView) - lipgloss.Height(shortcutsView) - 1).Render("")

	doc.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			renderedTitle,
			projectNameView,
			projectDescriptionView,
			spaceFiller,
			shortcutsView,
		),
	)
	return windowContainer.Render(doc.String())
}
