package mikrotik

import (
	"net/http"

	"github.com/ab22/stormrage/services/mikrotik"
)

type Handler interface {
	GetClients(w http.ResponseWriter, r *http.Request) error
}

// handler contains the websocket upgrader to handler mikrotik's websocket client.
type handler struct {
	mikrotikService mikrotik.Service
}

// NewHandler creates a new instance of Handler.
func NewHandler(s mikrotik.Service) Handler {
	return &handler{
		mikrotikService: s,
	}
}
