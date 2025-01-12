package main

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
