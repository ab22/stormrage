package auth

import (
	"net/http"

	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/services/auth"
)

type Handler interface {
	CheckAuth(w http.ResponseWriter, r *http.Request) error
	Login(w http.ResponseWriter, r *http.Request) error
	Logout(w http.ResponseWriter, r *http.Request) error
}

// handler contains all handlers in charge of authentication and sessions.
type handler struct {
	authService auth.Service
	cfg         *config.Config
}

// NewHandler creates a new Handler.
func NewHandler(authService auth.Service, cfg *config.Config) Handler {
	return &handler{
		authService: authService,
		cfg:         cfg,
	}
}
