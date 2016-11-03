package auth

import (
	"net/http"

	"github.com/ab22/stormrage/config"
)

type Handler interface {
	CheckAuth(w http.ResponseWriter, r *http.Request) error
}

// handler contains all handlers in charge of authentication and sessions.
type handler struct {
}

// NewHandler creates a new Handler.
func NewHandler(cfg *config.Config) Handler {
	return &handler{}
}
