package static

import (
	"fmt"
	"net/http"
)

// IndexHandler handles request for the index static file.
//
// Since Go's router sends all lost requests to home path '/',
// then we check if the URL path is not '/'.
// If the requested URL is '/', then we render the index.html template.
// If it's not, then we return a 404 response.
func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return nil
	}

	fmt.Fprintf(w, "<h1>Abemar</h1><br /><b>/</b>")

	return nil
}
