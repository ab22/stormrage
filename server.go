package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers"
	"github.com/ab22/stormrage/handlers/httputils"
	"github.com/ab22/stormrage/routes"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"
)

type Server struct {
	cfg         *config.Config
	router      *mux.Router
	cookieStore *sessions.CookieStore
	db          *gorm.DB
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

	log.Println("Configuring database...")
	err = server.createDatabaseConnection()
	if err = server.createDatabaseConnection(); err != nil {
		return nil, err
	}

	log.Println("Configuring router...")
	if err = server.configureRouter(); err != nil {
		return nil, err
	}

	log.Println("Creating static file server...")
	server.createStaticFilesServer()
	server.configureCookieStore()

	return server, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(
		fmt.Sprintf(":%d", s.cfg.Port),
		s.router,
	)
}

// createDatabaseConn creates a new GORM database with the specified database
// configuration.
func (s *Server) createDatabaseConnection() error {
	var (
		err              error
		dbCfg            = s.cfg.DB
		connectionString = fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Name,
		)
	)

	s.db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = s.db.DB().Ping()
	if err != nil {
		return err
	}

	s.db.DB().SetMaxIdleConns(10)
	s.db.LogMode(dbCfg.LogMode)

	return nil
}

func (s *Server) configureRouter() error {
	s.router = mux.NewRouter().StrictSlash(true)
	r, err := routes.NewRoutes(s.cfg, s.db)

	if err != nil {
		return err
	}

	s.bindRoutes(r.HTMLRoutes, false)
	s.bindRoutes(r.APIRoutes, true)

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
			ctx               = r.Context()
			commonMiddlewares = []handlers.MiddlewareFunc{
				handlers.HandleHTTPError,
				handlers.GzipContent,
			}
		)

		ctx = context.WithValue(ctx, "cookieStore", s.cookieStore)
		ctx = context.WithValue(ctx, "config", s.cfg)
		r = r.WithContext(ctx)

		for _, middleware := range commonMiddlewares {
			handler = middleware(handler)
		}

		if route.RequiresAuth() {
			handler = handlers.ValidateAuth(handler)
		}

		return handler(w, r)
	}
}

// createStaticFilesServer creates a static file server to serve all of the
// frontend files(html, js, css, etc).
func (s *Server) createStaticFilesServer() {
	var (
		staticFilesPath   = path.Join(s.cfg.FrontendAppPath, "static")
		commonMiddlewares = []handlers.MiddlewareFunc{
			handlers.HandleHTTPError,
			handlers.GzipContent,
			//handlers.NoDirListing,
		}
	)

	handler := func(w http.ResponseWriter, r *http.Request) error {
		file := path.Join(staticFilesPath, r.URL.Path)

		http.ServeFile(w, r, file)
		return nil
	}

	for _, middleware := range commonMiddlewares {
		handler = middleware(handler)
	}

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "cookieStore", s.cookieStore)
		ctx = context.WithValue(ctx, "config", s.cfg)
		r = r.WithContext(ctx)

		err := handler(w, r)

		if err != nil {
			log.Printf("static file handler [%s][%s] returned error: %s", r.Method, r.URL.Path, err)
			httputils.WriteError(w, http.StatusInternalServerError, "")
		}
	})

	s.router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static", httpHandler))
}

// configureCookieStore creates the cookie store used to validate user
// sessions.
func (s *Server) configureCookieStore() {
	secretKey := s.cfg.Secret

	gob.Register(&handlers.SessionData{})

	s.cookieStore = sessions.NewCookieStore([]byte(secretKey))
	s.cookieStore.MaxAge(0)
}
