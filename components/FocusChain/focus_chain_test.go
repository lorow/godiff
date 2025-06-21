package FocusChain

import "testing"

type FocusableTestItem struct {
	isFocused bool
}

func (f *FocusableTestItem) Focus() {
	f.isFocused = true
}

func (f *FocusableTestItem) Blur() {
	f.isFocused = false
}

func (f FocusableTestItem) HasFocusChain() bool {
	return false
}

func (f FocusableTestItem) HandleNext() *ReachedFocusChainLimit {
	return nil
}

func (f FocusableTestItem) HandlePrevious() *ReachedFocusChainLimit {
	return nil
}

func TestMovingFocus(t *testing.T) {

	firstItem := &FocusableTestItem{isFocused: true}
	secondItem := &FocusableTestItem{isFocused: false}

	focusChain := New(WithItem(firstItem), WithItem(secondItem))

	if focusChain.GetCurrentlyFocused() != firstItem {
		t.Errorf("Expected firstItem to be focused, but got %v", focusChain.GetCurrentlyFocused())
	}

	if !firstItem.isFocused {
		t.Errorf("Expected firstItem to be focused, but got %v", firstItem.isFocused)
	}

	focusChain.Next()

	if focusChain.GetCurrentlyFocused() != secondItem {
		t.Errorf("Expected firstItem to be focused, but got %v", focusChain.GetCurrentlyFocused())
	}

	if !secondItem.isFocused {
		t.Errorf("Expected secondItem to be focused, but got %v", secondItem.isFocused)
	}

}
