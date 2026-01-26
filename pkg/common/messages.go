package Shared

type SetNewSizeMsg struct {
	Width, Height int
}

type RouteMsg struct {
	To      Screen
	Payload any
}
