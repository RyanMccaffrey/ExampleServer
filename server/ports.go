package server

import "net/http"

type Route struct {
	Path    string
	Method  string
	Wrapper *HandlerWrapper // Optional wrapper for the handler
	Handler http.HandlerFunc
}

type HandlerWrapper func(handler http.HandlerFunc) http.HandlerFunc

type WebService interface {
	ServeRoutes() error
}

type RemoteRoutes interface {
	AddRouteHandler(path string, method string, wrapper *HandlerWrapper, handler http.HandlerFunc) error
	AddRoutes(routes []Route) error
	HandleRoutes() error
}

type RemoteServer interface {
	Serve(addr string) error
}
