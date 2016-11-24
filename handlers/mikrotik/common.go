package mikrotik

import (
	"net/http"

	"github.com/ab22/stormrage/handlers/httputils"
)

func (h *handler) GetClients(w http.ResponseWriter, r *http.Request) error {
	clients, err := h.mikrotikService.RequestClients()

	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, clients)
}
