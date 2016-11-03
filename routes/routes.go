package routes

import (
	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers/static"
)

// Routes contains all HTML and API routes for the application.
type Routes struct {
	HTMLRoutes []Route
	// APIRoutes  []Route
}

// NewRoutes creates a new Router instance and initializes all HTML
// and API Routes.
func NewRoutes(_ *config.Config) (*Routes, error) {
	var (
		staticHandler = static.NewHandler()
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
	}, nil
}
