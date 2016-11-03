package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers"
	"github.com/ab22/stormrage/handlers/httputils"
	"github.com/ab22/stormrage/routes"
	"github.com/gorilla/mux"
)

type Server struct {
	cfg    *config.Config
	router *mux.Router
}

func NewServer() (*Server, error) {
	var (
		err    error
		server = &Server{}
	)

	log.Println("Configuring server...")
	server.cfg, err = config.New()
	server.cfg.Print()

	if err != nil {
		return nil, err
	}

	err = server.configureRouter()

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(
		fmt.Sprintf(":%d", s.cfg.Port),
		s.router,
	)
}

func (s *Server) configureRouter() error {
	s.router = mux.NewRouter().StrictSlash(true)
	r, err := routes.NewRoutes(s.cfg)

	if err != nil {
		return err
	}

	s.bindRoutes(r.HTMLRoutes, false)

	return nil
}

// bindRoutes adds all routes to the server's router.
func (s *Server) bindRoutes(r []routes.Route, apiRoute bool) {
	for _, route := range r {
		var routePath string
		httpHandler := s.makeHTTPHandler(route)

		if apiRoute {
			routePath = "/api/" + route.Pattern()
		} else {
			routePath = route.Pattern()
		}

		s.router.
			Methods(route.Method()).
			Path(routePath).
			HandlerFunc(httpHandler)
	}
}

// makeHTTPHandler creates a http.HandlerFunc from a custom http function and logs the error if
// exists: func(http.ResponseWriter, *http.Request) error.
func (s *Server) makeHTTPHandler(route routes.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc := s.handleWithMiddlewares(route)
		err := handlerFunc(w, r)

		if err != nil {
			log.Printf("Handler [%s][%s] returned error: %s", r.Method, r.URL.Path, err)
		}
	}
}

// handleWithMiddlewares applies all middlewares to the specified route. Some
// middleware functions are applied depending on the route's properties, such
// as ValidateAuth and Authorize middlewares. These last 2 functions require
// that the route RequiresAuth() and that RequiredRoles() > 0.
func (s *Server) handleWithMiddlewares(route routes.Route) httputils.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			handler           = route.HandlerFunc()
			commonMiddlewares = []handlers.MiddlewareFunc{
				handlers.HandleHTTPError,
				handlers.GzipContent,
			}
		)

		for _, middleware := range commonMiddlewares {
			handler = middleware(handler)
		}

		return handler(w, r)
	}
}
