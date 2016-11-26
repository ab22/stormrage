package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers/httputils"
	"github.com/gorilla/sessions"
)

// MiddlewareFunc describes a function that takes a ContextHandler and
// returns a ContextHandler.
//
// The idea of a middleware function is to validate/read/modify data before or
// after calling the next middleware function.
type MiddlewareFunc func(httputils.HandlerFunc) httputils.HandlerFunc

// extendSessionLifetime determines if the session's lifetime needs to be
// extended. Session's lifetime should be extended only if the session's
// current lifetime is below sessionLifeTime/2. Returns true if the session
// needs to be extended.
func extendSessionLifetime(sessionData *SessionData, sessionLifeTime time.Duration) bool {
	return sessionData.ExpiresAt.Sub(time.Now()) <= sessionLifeTime/2
}

// ValidateAuth validates that the user cookie is set up before calling the
// handler passed as parameter.
func ValidateAuth(h httputils.HandlerFunc) httputils.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			ctx     = r.Context()
			cfg, ok = ctx.Value("config").(*config.Config)
		)

		if !ok {
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return fmt.Errorf("validate auth: error casting config object")
		}

		cookieStore, ok := ctx.Value("cookieStore").(*sessions.CookieStore)

		if !ok {
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return fmt.Errorf("validate auth: could not cast value as cookie store: %s", ctx.Value("cookieStore"))
		}

		session, err := cookieStore.Get(r, cfg.SessionCookieName)

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		}

		sessionData, ok := session.Values["data"].(*SessionData)

		if !ok || sessionData.IsInvalid() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		} else if time.Now().After(sessionData.ExpiresAt) {
			session.Options.MaxAge = -1
			session.Save(r, w)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return nil
		}

		// Save session only if the session was extended.
		if extendSessionLifetime(sessionData, cfg.SessionLifeTime) {
			sessionData.ExpiresAt = time.Now().Add(cfg.SessionLifeTime)
			session.Save(r, w)
		}

		ctx = context.WithValue(ctx, "sessionData", sessionData)
		authenticatedRequest := r.WithContext(ctx)
		return h(w, authenticatedRequest)
	}
}

// HandleHTTPError sets the appropriate headers to the response if a http
// handler returned an error. This might be used in the future if different
// types of errors are returned.
func HandleHTTPError(h httputils.HandlerFunc) httputils.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := h(w, r)

		if err != nil {
			httputils.WriteError(w, http.StatusInternalServerError, "")
		}

		return err
	}
}
