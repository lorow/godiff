package FocusChain

import tea "github.com/charmbracelet/bubbletea"

type Focusable interface {
	Focus() tea.Cmd
	Blur() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
}

type ReachedFocusChainEnd struct{}

type Model struct {
	index int
	chain []Focusable
}

func New(ops ...func(*Model)) *Model {
	chain := &Model{index: 0, chain: []Focusable{}}

	for _, op := range ops {
		op(chain)
	}

	return chain
}

func WithItem(item Focusable) func(chain *Model) {
	return func(chain *Model) {
		chain.chain = append(chain.chain, item)
	}
}

func (chain *Model) GetCurrentlySelected() Focusable {
	return chain.chain[chain.index]
}

func (chain *Model) Next() (tea.Cmd, *ReachedFocusChainEnd) {

	if len(chain.chain) == 0 {
		return nil, &ReachedFocusChainEnd{}
	}

	if chain.index < len(chain.chain)-1 {
		cmds := []tea.Cmd{}
		cmds = append(cmds, chain.chain[chain.index].Blur())
		chain.index++
		cmds = append(cmds, chain.chain[chain.index].Focus())
		return tea.Batch(cmds...), nil
	}

	return nil, &ReachedFocusChainEnd{}
}

func (chain *Model) Previous() (tea.Cmd, *ReachedFocusChainEnd) {

	if len(chain.chain) == 0 {
		return nil, &ReachedFocusChainEnd{}
	}

	if chain.index > 0 {
		cmds := []tea.Cmd{}
		cmds = append(cmds, chain.chain[chain.index].Blur())
		chain.index--
		cmds = append(cmds, chain.chain[chain.index].Focus())
		return tea.Batch(cmds...), nil
	}

	return nil, &ReachedFocusChainEnd{}

}

func (chain *Model) JumpFocus(index string) (tea.Cmd, *ReachedFocusChainEnd) {
	return nil, nil
}
