package FocusChain

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type FocusableTestItem struct {
	isFocused bool
}

func (f *FocusableTestItem) Focus() {
	f.isFocused = true
}

func (f *FocusableTestItem) Blur() {
	f.isFocused = false
}

func (f *FocusableTestItem) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func TestMovingFocusForward(t *testing.T) {

	firstItem := &FocusableTestItem{isFocused: true}
	secondItem := &FocusableTestItem{isFocused: false}

	focusChain := New(WithItem(firstItem), WithItem(secondItem))

	if focusChain.GetCurrentlySelected() != firstItem {
		t.Errorf("Expected firstItem to be focused, but got %v", focusChain.GetCurrentlySelected())
	}

	if !firstItem.isFocused {
		t.Errorf("Expected firstItem to be focused, but got %v", firstItem.isFocused)
	}

	focusChain.Next()

	if focusChain.GetCurrentlySelected() != secondItem {
		t.Errorf("Expected secondItem to be focused, but got %v", focusChain.GetCurrentlySelected())
	}

	if !secondItem.isFocused {
		t.Errorf("Expected secondItem to be focused, but got %v", secondItem.isFocused)
	}
}

func TestMovingFocusBackwards(t *testing.T) {

	firstItem := &FocusableTestItem{isFocused: true}
	secondItem := &FocusableTestItem{isFocused: false}

	focusChain := New(WithItem(firstItem), WithItem(secondItem))

	if focusChain.GetCurrentlySelected() != firstItem {
		t.Errorf("Expected firstItem to be currently selected, but got %v", focusChain.GetCurrentlySelected())
	}

	if !firstItem.isFocused {
		t.Errorf("Expected firstItem to be focused, but got %v", firstItem.isFocused)
	}

	focusChain.Next()

	if focusChain.GetCurrentlySelected() != secondItem {
		t.Errorf("Expected secondItem to be focused, but got %v", focusChain.GetCurrentlySelected())
	}

	if !secondItem.isFocused {
		t.Errorf("Expected secondItem to be focused, but got %v", secondItem.isFocused)
	}

	focusChain.Previous()

	if focusChain.GetCurrentlySelected() != firstItem {
		t.Errorf("Expected firstItem to be focused, but got %v", focusChain.GetCurrentlySelected())
	}

	if !firstItem.isFocused {
		t.Errorf("Expected firstItem to be focused, but got %v", secondItem.isFocused)
	}

}

func TestMovingFocusHitEnd(t *testing.T) {

	firstItem := &FocusableTestItem{isFocused: true}
	focusChain := New(WithItem(firstItem))
	result := focusChain.Next()

	if result == nil {
		t.Errorf("Expected result to indicated the end of focus chain, but got %v", result)
	}

	// generally, once we hit the end, the component should not be blurred just yet
	// that should be handled by either the higher focus chain, or the view
	if !firstItem.isFocused {
		t.Errorf("Expected result to indicated the end of focus chain, but got %v", result)
	}

	result = focusChain.Previous()

	if result == nil {
		t.Errorf("Expected result to indicated the end of focus chain, but got %v", result)
	}

	// same goes for the Previous() mechanic
	if !firstItem.isFocused {
		t.Errorf("Expected result to indicated the end of focus chain, but got %v", result)
	}

}
