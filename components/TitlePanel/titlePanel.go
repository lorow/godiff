package TitlePanel

type Model struct {
	styles Styles
	title  string
}

func New(ops ...func(*Model)) *Model {
	model := &Model{
		styles: DefaultStyles(),
		title:  "",
	}

	for _, op := range ops {
		op(model)
	}

	return model
}

func WithTitle(title string) func(*Model) {
	return func(model *Model) {
		model.title = title
	}
}

func WithStyles(styles Styles) func(*Model) {
	return func(model *Model) {
		model.styles = styles
	}
}

func (m Model) View() string {
	title := m.styles.Title.Render(m.title)
	return m.styles.Container.Render(title)
}
