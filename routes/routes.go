package routes

import (
	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers/auth"
	"github.com/ab22/stormrage/handlers/mikrotik"
	"github.com/jinzhu/gorm"

	authservices "github.com/ab22/stormrage/services/auth"
	mikrotikservices "github.com/ab22/stormrage/services/mikrotik"
	userservices "github.com/ab22/stormrage/services/user"
	"github.com/ab22/stormrage/services/ws"
)

// NewRoutes creates a new Router instance and initializes all API Routes.
func NewRoutes(cfg *config.Config, db *gorm.DB) ([]Route, error) {
	var (
		userService      = userservices.NewService(db)
		authService      = authservices.NewService(db, userService)
		mikrotikService  = mikrotikservices.NewService(cfg)
		websocketService = ws.NewServer()

		// staticHandler   = static.NewHandler(cfg)
		authHandler     = auth.NewHandler(authService, cfg)
		mikrotikHandler = mikrotik.NewHandler(mikrotikService)
	)

	// API routes
	return []Route{
		&route{
			pattern:      "ws/onConnect/",
			method:       "GET",
			handlerFunc:  websocketService.OnConnect,
			requiresAuth: true,
			gzipContent:  false,
		},
		&route{
			pattern:      "auth/checkAuthentication/",
			method:       "POST",
			handlerFunc:  authHandler.CheckAuth,
			requiresAuth: true,
			gzipContent:  true,
		},
		&route{
			pattern:      "auth/login/",
			method:       "POST",
			handlerFunc:  authHandler.Login,
			requiresAuth: false,
			gzipContent:  true,
		},
		&route{
			pattern:      "auth/logout/",
			method:       "POST",
			handlerFunc:  authHandler.Logout,
			requiresAuth: false,
			gzipContent:  true,
		},
		&route{
			pattern:      "mikrotik/getClients/",
			method:       "POST",
			handlerFunc:  mikrotikHandler.GetClients,
			requiresAuth: true,
			gzipContent:  true,
		},
	}, nil
}
