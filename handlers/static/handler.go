package static

import "net/http"

type Handler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request) error
}

// handler contains all handlers in charge of serving static pages and files.
type handler struct {
}

// NewHandler creates a new instance of Handler.
func NewHandler() Handler {
	return &handler{}
}
