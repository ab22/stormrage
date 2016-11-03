package static

import (
	"html/template"
	"net/http"

	"github.com/ab22/stormrage/config"
)

type Handler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request) error
}

// handler contains all handlers in charge of serving static pages and files.
type handler struct {
	cachedTemplates *template.Template
}

// NewHandler creates a new instance of Handler.
func NewHandler(cfg *config.Config) Handler {
	frontendAppPath := cfg.FrontendAppPath

	return &handler{
		cachedTemplates: template.Must(template.ParseFiles(frontendAppPath + "/index.html")),
	}
}
