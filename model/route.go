package model

import (
	"net"
)

type Route struct {
	Network  *net.IPNet
	RouterID string
	Metric   int
}
type RouteList struct {
	Routes []*Route
}

func NewRouteList() *RouteList {
	return &RouteList{
		Routes: make([]*Route, 0),
	}
}
