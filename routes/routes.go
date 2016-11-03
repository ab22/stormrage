package routes

import (
	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers/auth"
	"github.com/ab22/stormrage/handlers/static"
	"github.com/jinzhu/gorm"

	authservices "github.com/ab22/stormrage/services/auth"
	userservices "github.com/ab22/stormrage/services/user"
)

// Routes contains all HTML and API routes for the application.
type Routes struct {
	HTMLRoutes []Route
	APIRoutes  []Route
}

// NewRoutes creates a new Router instance and initializes all HTML
// and API Routes.
func NewRoutes(cfg *config.Config, db *gorm.DB) (*Routes, error) {
	var (
		userService = userservices.NewService(db)
		authService = authservices.NewService(db, userService)

		staticHandler = static.NewHandler(cfg)
		authHandler   = auth.NewHandler(authService, cfg)
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
			&route{
				pattern:      "auth/login/",
				method:       "POST",
				handlerFunc:  authHandler.Login,
				requiresAuth: false,
			},
			&route{
				pattern:      "auth/logout/",
				method:       "POST",
				handlerFunc:  authHandler.Logout,
				requiresAuth: false,
			},
		},
	}, nil
}
