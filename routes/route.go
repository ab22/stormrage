package routes

import "net/http"

// Route interface.
type Route interface {
	Pattern() string
	Method() string
	HandlerFunc() func(http.ResponseWriter, *http.Request) error
	RequiresAuth() bool
	RequiredRoles() []string
}

type route struct {
	pattern       string
	method        string
	handlerFunc   func(http.ResponseWriter, *http.Request) error
	requiresAuth  bool
	requiredRoles []string
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

func (r *route) RequiredRoles() []string {
	return r.requiredRoles
}
