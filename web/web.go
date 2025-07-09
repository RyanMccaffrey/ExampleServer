package web

import (
	"net/http"

	"github.com/RyanMccaffrey/ExampleServer/server"
)

type HelloWorldService struct {
	addressPrefix string
}

func NewHelloWorldService(addressPrefix string) *HelloWorldService {
	return &HelloWorldService{
		addressPrefix: addressPrefix,
	}
}

func (s *HelloWorldService) ServeRoutes() error {
	routeService := server.NewExampleRouteService(s.addressPrefix, nil)
	if err := routeService.AddRouteHandler("/", "GET", nil, func(w http.ResponseWriter, r *http.Request) {
		response := "Hello, World!"
		w.Write([]byte(response))
	}); err != nil {
		return err
	}
	if err := routeService.HandleRoutes(); err != nil {
		return err
	}
	return nil
}
