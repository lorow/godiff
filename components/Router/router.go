package Router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type UpdateMessage tea.Msg

type RouteMsg struct {
	route      string
	updMessage UpdateMessage
}

func RouteTo(route string, updateMessage UpdateMessage) tea.Cmd {
	return func() tea.Msg {
		return RouteMsg{
			route, updateMessage,
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

func (r *Model) HandleRouteTo(msg RouteMsg) tea.Cmd {
	if _, ok := r.views[msg.route]; ok {
		log.Debug("Routing to", "route", msg.route, "update msg", msg.updMessage)
		r.currentRoute = msg.route
		view := r.GetCurrentVIew()
		view, cmd := view.Update(msg.updMessage)
		r.UpdateRoute(msg.route, view)

		return cmd
	}

	log.Debug("Route not found:", "route", msg.route)
	return nil
}
