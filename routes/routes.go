package routes

import (
	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers/auth"
	"github.com/ab22/stormrage/handlers/mikrotik"
	"github.com/ab22/stormrage/handlers/static"
	"github.com/jinzhu/gorm"

	authservices "github.com/ab22/stormrage/services/auth"
	mikrotikservices "github.com/ab22/stormrage/services/mikrotik"
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
		userService     = userservices.NewService(db)
		authService     = authservices.NewService(db, userService)
		mikrotikService = mikrotikservices.NewService(cfg)

		staticHandler   = static.NewHandler(cfg)
		authHandler     = auth.NewHandler(authService, cfg)
		mikrotikHandler = mikrotik.NewHandler(mikrotikService)
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
			&route{
				pattern:      "mikrotik/getClients/",
				method:       "POST",
				handlerFunc:  mikrotikHandler.GetClients,
				requiresAuth: true,
			},
		},
	}, nil
}
