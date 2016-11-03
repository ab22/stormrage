package httputils

import (
	"encoding/json"
	"io"
	"net/http"
)

// A RouteType is used to differentiate requests to the API from the statics
// file server.
type RouteType int

const (
	// API type call.
	API RouteType = iota
	// Static type call.
	Static
)

// HandlerFunc defines how API handler functions should be defined. ApiHandler
// functions should return any kind of value which will be turned into json
// and an *ApiError.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// WriteError writes an error to the ResponseWriter. If no message is
// specified, then we retrieve the default status text from the specified
// code parameter.
func WriteError(w http.ResponseWriter, code int, errMsg string) {
	if errMsg == "" {
		errMsg = http.StatusText(code)
	}

	http.Error(w, errMsg, code)
}

// WriteJSON writes the specified data as json.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

// DecodeJSON is a handy method to decode json into a interface{}.
func DecodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
