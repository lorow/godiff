package FocusChain

import tea "github.com/charmbracelet/bubbletea"

type Focusable interface {
	Focus()
	Blur()
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

func (chain *Model) Next() *ReachedFocusChainEnd {

	if len(chain.chain) == 0 {
		return &ReachedFocusChainEnd{}
	}

	if chain.index < len(chain.chain)-1 {
		chain.chain[chain.index].Blur()
		chain.index++
		chain.chain[chain.index].Focus()
		return nil
	}

	return &ReachedFocusChainEnd{}
}

func (chain *Model) Previous() *ReachedFocusChainEnd {

	if len(chain.chain) == 0 {
		return &ReachedFocusChainEnd{}
	}

	if chain.index > 0 {
		chain.chain[chain.index].Blur()
		chain.index--
		chain.chain[chain.index].Focus()
		return nil
	}

	return &ReachedFocusChainEnd{}

}

func (chain *Model) JumpFocus(index string) *ReachedFocusChainEnd {
	return nil
}
