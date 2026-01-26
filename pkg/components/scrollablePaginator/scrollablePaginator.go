package scrollablePaginator

type Model struct {
	viewWindowBounds [2]int
}

func New() Model {
	return Model{}
}

func (m Model) GetViewWindowBounds() [2]int {
	return m.viewWindowBounds
}

func (m Model) GetLowerBound() int {
	return m.viewWindowBounds[0]
}

func (m Model) GetUpperBound() int {
	return m.viewWindowBounds[1]
}

func (m *Model) MoveViewWindowDown() {
	m.viewWindowBounds[0]++
	m.viewWindowBounds[1]++
}

func (m *Model) MoveViewWindowUp() {
	m.viewWindowBounds[0]--
	m.viewWindowBounds[1]--
}

func (m *Model) SetViewWindowBounds(bounds [2]int) {
	m.viewWindowBounds = bounds
}

func (m *Model) RecalculateViewBounds(availableHeight, itemHeight int) {
	// we need to update the max visible area

	distanceBetweenCurrentBounds := m.viewWindowBounds[1] - m.viewWindowBounds[0]
	newMaxDisntance := availableHeight / itemHeight

	// now, if we're just starting out, we only need to set the upper bound
	if m.viewWindowBounds[0] == 0 && m.viewWindowBounds[1] == 0 {
		m.viewWindowBounds[1] = max(0, newMaxDisntance)
		return
	}

	// otherwise, we're updating the list, downwards.
	expandBy := newMaxDisntance - distanceBetweenCurrentBounds
	m.viewWindowBounds[1] += expandBy
}
