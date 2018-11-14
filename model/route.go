package model;

type Route struct{}
type RouteList struct {
	Routes []Route
}

func NewRouteList() *RouteList {
	return &RouteList{
		Routes: make([]Route, 0),
	}
}
