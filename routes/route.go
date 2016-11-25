package routes

import "net/http"

// Route interface.
type Route interface {
	Pattern() string
	Method() string
	HandlerFunc() func(http.ResponseWriter, *http.Request) error
	RequiresAuth() bool
	GzipContent() bool
}

type route struct {
	pattern      string
	method       string
	handlerFunc  func(http.ResponseWriter, *http.Request) error
	requiresAuth bool
	gzipContent  bool
}

func (r *route) Pattern() string {
	return r.pattern
}

func (r *route) Method() string {
	return r.method
}

func (r *route) HandlerFunc() func(http.ResponseWriter, *http.Request) error {
	return r.handlerFunc
}

func (r *route) RequiresAuth() bool {
	return r.requiresAuth
}

func (r *route) GzipContent() bool {
	return r.gzipContent
}
