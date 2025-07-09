package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	pathRegex   = regexp.MustCompile(`^\/([A-Za-z0-9\-\._~!$&'()*+,;=:@/]*)?$`)
	httpMethods = map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodHead:    true,
		http.MethodOptions: true,
		http.MethodConnect: true,
		http.MethodTrace:   true,
	}
)

func validatePathAndMethod(path, method string) error {
	if !pathRegex.MatchString(path) {
		return errors.New(fmt.Sprintf("invalid HTTP path format: %s", path))
	}

	if !httpMethods[strings.ToUpper(method)] {
		return errors.New(fmt.Sprintf("invalid HTTP method: %s", method))
	}
	return nil
}

type ExampleRouteService struct {
	addrPrefix string
	routes     []Route
	wrapper    *HandlerWrapper
}

func NewExampleRouteService(addrPrefix string, wrapper *HandlerWrapper) *ExampleRouteService {
	return &ExampleRouteService{
		addrPrefix: addrPrefix,
		routes:     []Route{},
		wrapper:    wrapper,
	}
}
func (s *ExampleRouteService) AddRouteHandler(path string, method string, wrapper *HandlerWrapper, handler http.HandlerFunc) error {
	if err := validatePathAndMethod(path, method); err != nil {
		return err
	}
	s.routes = append(s.routes, Route{
		Path:    path,
		Method:  method,
		Wrapper: wrapper,
		Handler: handler,
	})
	return nil
}
func (s *ExampleRouteService) AddRoutes(routes []Route) error {
	if len(routes) == 0 {
		return errors.New("no routes provided")
	}
	for _, route := range routes {
		if err := validatePathAndMethod(route.Path, route.Method); err != nil {
			return err
		}
		s.routes = append(s.routes, route)
	}
	return nil
}
func (s *ExampleRouteService) HandleRoutes() error {
	if len(s.routes) == 0 {
		return errors.New("no routes to handle")
	}
	consolidatedRoutes, err := s.consolidateRoutes()
	if err != nil {
		return err
	}

	for _, routes := range consolidatedRoutes {
		s.handleRouteMethods(routes)
	}
	return nil
}

func (s *ExampleRouteService) consolidateRoutes() (map[string][]Route, error) {
	routeMethods := make(map[string][]Route)
	for _, route := range s.routes {
		if existing, exists := routeMethods[route.Path]; exists {
			for _, r := range existing {
				if r.Method == route.Method {
					return nil, fmt.Errorf("duplicate route found: %s %s", route.Method, route.Path)
				}
			}
			routeMethods[route.Path] = append(routeMethods[route.Path], route)
		} else {
			routeMethods[route.Path] = []Route{route}
		}
	}
	return routeMethods, nil
}

func (s *ExampleRouteService) handleRouteMethods(routes []Route) {

	// all routes should be identical in path, so we can just use the first one
	http.HandleFunc(s.addrPrefix+routes[0].Path, func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if r.Method == route.Method {
				handler := s.wrapHandler(route.Wrapper, route.Handler)
				handler(w, r)
				return
			}
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})
}
func (s *ExampleRouteService) wrapHandler(wrapper *HandlerWrapper, handler http.HandlerFunc) http.HandlerFunc {
	if s.wrapper != nil {
		handler = (*s.wrapper)(handler)
	}
	if wrapper != nil {
		return (*wrapper)(handler)
	}
	return handler
}
