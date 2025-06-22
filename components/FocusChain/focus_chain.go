package FocusChain

type Focusable interface {
	Focus()
	Blur()
	HasFocusChain() bool
	HandleNext() *ReachedFocusChainLimit
	HandlePrevious() *ReachedFocusChainLimit
}

type ReachedFocusChainLimit struct{}

// FocusChain stores
type FocusChain struct {
	index int
	chain []Focusable
}

func New(ops ...func(*FocusChain)) *FocusChain {
	chain := &FocusChain{index: 0, chain: []Focusable{}}

	for _, op := range ops {
		op(chain)
	}

	return chain
}

func WithItem(item Focusable) func(chain *FocusChain) {
	return func(chain *FocusChain) {
		chain.chain = append(chain.chain, item)
	}
}

func (chain *FocusChain) GetCurrentlySelected() Focusable {
	return chain.chain[chain.index]
}

func (chain *FocusChain) Next() *ReachedFocusChainLimit {

	if len(chain.chain) == 0 {
		return &ReachedFocusChainLimit{}
	}

	currentItem := chain.chain[chain.index]

	if currentItem.HasFocusChain() {
		result := currentItem.HandleNext()
		if result != nil {
			currentItem.Blur()
		}
	}

	if chain.index < len(chain.chain)-1 {
		chain.index++
		chain.chain[chain.index].Focus()
		return nil
	}

	return &ReachedFocusChainLimit{}
}

func (chain *FocusChain) Previous() *ReachedFocusChainLimit {

	if len(chain.chain) == 0 {
		return &ReachedFocusChainLimit{}
	}

	currentItem := chain.chain[chain.index]

	if currentItem.HasFocusChain() {
		result := currentItem.HandlePrevious()
		if result != nil {
			currentItem.Blur()
		}
	}

	if chain.index > 0 {
		chain.index--
		chain.chain[chain.index].Focus()
		return nil
	}

	return &ReachedFocusChainLimit{}

}

func (chain *FocusChain) JumpFocus(index string) *ReachedFocusChainLimit {
	return nil
}
