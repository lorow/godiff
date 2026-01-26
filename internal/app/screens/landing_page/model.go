package landingPage

import (
	"godiff/internal/db"
	Shared "godiff/pkg/common"
	"godiff/pkg/components/FocusChain"
	"godiff/pkg/components/ItemList"
	"godiff/pkg/components/ShortcutsPanel"
	"godiff/pkg/components/TextInput"
	"godiff/pkg/components/TitlePanel"
)

type LandingPageState int

const (
	Initial LandingPageState = iota
	LoadingProjects
	LoadedProjects
	ErrorLoadingProjects
)

type Model struct {
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

func New() Model {
	searchInput := TextInput.New(TextInput.WIthOnSubmit(onSearchSubmit))
	searchInput.Focus()

	itemList := ItemList.New(
		ItemList.WithTittle("Projects"),
		ItemList.WithNoItemsText("Empty"),
		ItemList.WithOnSelection(onItemListSelect),
	)

	titlePanel := TitlePanel.New(TitlePanel.WithTitle("Welcome to GoDiff - 1.0.0"))

	basicShortcuts := []ShortcutsPanel.Shortcut{
		ShortcutsPanel.NewShortcut("↑/↓", Shared.NoopShortcut, "Change focus", nil), // display only
		ShortcutsPanel.NewShortcut("^Q", Shared.ExitShortcut, "Quit", Shared.ExitCmd()),
		ShortcutsPanel.NewShortcut("^N", Shared.NewProjectShortcut, "Create new project", Shared.RouteTo(Shared.NewProjectWizardScreen, nil)),
		ShortcutsPanel.NewShortcut("^O", Shared.JumpFocusShortcut, "Jump focus", Shared.JumpFocusCmd()),
	}

	onProjectSelectShortcuts := []ShortcutsPanel.Shortcut{
		ShortcutsPanel.NewShortcut("↑/↓", Shared.NoopShortcut, "Change focus", nil), // display only
		ShortcutsPanel.NewShortcut("^Q", Shared.ExitShortcut, "Quit", Shared.ExitCmd()),
		ShortcutsPanel.NewShortcut("^N", Shared.NewProjectShortcut, "Create new project", Shared.RouteTo(Shared.NewProjectWizardScreen, nil)),
		ShortcutsPanel.NewShortcut("Enter", Shared.NoopShortcut, "Launch selected project", nil), // display only
		ShortcutsPanel.NewShortcut("^O", Shared.JumpFocusShortcut, "Jump focus", Shared.JumpFocusCmd()),
	}

	shortcutsPanel := ShortcutsPanel.New(
		ShortcutsPanel.WithShortcuts(onProjectSelectShortcuts),
	)

	focusChain := FocusChain.New(FocusChain.WithItem(searchInput), FocusChain.WithItem(itemList))

	return Model{
		searchInput:              searchInput,
		itemList:                 itemList,
		titlePanel:               titlePanel,
		shortcutsPanel:           shortcutsPanel,
		basicShortcuts:           basicShortcuts,
		onProjectSelectShortcuts: onProjectSelectShortcuts,
		focusChain:               focusChain,
	}
}
