package newProjectWizard

import (
	Shared "godiff/pkg/common"
	"godiff/pkg/components/FocusChain"
	"godiff/pkg/components/ShortcutsPanel"
	"godiff/pkg/components/TextInput"
	"godiff/pkg/components/TitlePanel"
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
		ShortcutsPanel.NewShortcut("esc", Shared.EscapeShortcut, "Go back", Shared.RouteTo(Shared.LandingScreen, nil)),
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
