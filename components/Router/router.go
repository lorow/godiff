package Router

import tea "github.com/charmbracelet/bubbletea"

type additionalData interface{}

type RouteMsg struct {
	route string
	data  additionalData
}

func RouteTo(route string, additionalData additionalData) tea.Cmd {
	return func() tea.Msg {
		return RouteMsg{
			route, additionalData,
		}
	}
}

type Model struct {
	currentRoute string
	views        map[string]tea.Model
}

func New(options ...func(*Model)) *Model {
	router := &Model{
		currentRoute: "/",
		views:        make(map[string]tea.Model),
	}

	for _, opt := range options {
		opt(router)
	}

	return router
}

func WithStartingPage(view tea.Model) func(*Model) {
	return func(r *Model) {
		r.views["/"] = view
	}
}

func WithRegisterRoute(route string, view tea.Model) func(*Model) {
	return func(r *Model) {
		r.views[route] = view
	}
}

func (r *Model) GetCurrentVIew() tea.Model {
	return r.views[r.currentRoute]
}

func (r *Model) GetCUrrentRoute() string {
	return r.currentRoute
}

// to be used mainly for mass initialization of the models
func (r *Model) GetViews() map[string]tea.Model {
	return r.views
}

func (r *Model) UpdateRoute(route string, view tea.Model) {
	r.views[route] = view
}

func (r *Model) HandleRouteTo(msg tea.Msg) {
	// todo add handling of route to message
}
