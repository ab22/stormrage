package routes

import (
	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers/auth"
	"github.com/ab22/stormrage/handlers/static"
)

// Routes contains all HTML and API routes for the application.
type Routes struct {
	HTMLRoutes []Route
	APIRoutes  []Route
}

// NewRoutes creates a new Router instance and initializes all HTML
// and API Routes.
func NewRoutes(cfg *config.Config) (*Routes, error) {
	var (
		staticHandler = static.NewHandler(cfg)
		authHandler   = auth.NewHandler(cfg)
	)

	return &Routes{
		HTMLRoutes: []Route{
			&route{
				pattern:      "/",
				method:       "GET",
				handlerFunc:  staticHandler.IndexHandler,
				requiresAuth: false,
			},
		},
		APIRoutes: []Route{
			&route{
				pattern:      "auth/checkAuthentication/",
				method:       "POST",
				handlerFunc:  authHandler.CheckAuth,
				requiresAuth: true,
			},
		},
	}, nil
}
