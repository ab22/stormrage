package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/ab2_2/playcheck/config"
	"bitbucket.org/ab2_2/playcheck/handlers/httputils"
	"github.com/gorilla/sessions"
)

// MiddlewareFunc describes a function that takes a ContextHandler and
// returns a ContextHandler.
//
// The idea of a middleware function is to validate/read/modify data before or
// after calling the next middleware function.
type MiddlewareFunc func(httputils.HandlerFunc) httputils.HandlerFunc

// NoDirListing helps avoid listing folder directories.
//
// Go's http.FileServer by default, lists the directories and files
// of the specified folder to serve and cannot be disabled.
// To prevent directory listing, noDirListing checks if the
// path requests ends in '/'. If it does, then the client is requesting
// to explore a folder and we return a 404 (Not found), else, we just
// call the http.Handler passed as parameter.
func NoDirListing(h httputils.HandlerFunc) httputils.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		urlPath := r.URL.Path

		if urlPath == "" || strings.HasSuffix(urlPath, "/") {
			http.NotFound(w, r)
			return nil
		}

		return h(w, r)
	}
}

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

// GzipContent is a middleware function for handlers to encode content to gzip.
func GzipContent(h httputils.HandlerFunc) httputils.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Add("Vary", "Accept-Encoding")

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			return h(w, r)
		}

		w.Header().Set("Content-Encoding", "gzip")

		gzipResponseWriter := httputils.NewGzipResponseWriter(w)
		defer gzipResponseWriter.Close()

		return h(gzipResponseWriter, r)
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

// ForceSSL forces HTTPS on the Heroku production servers.
//
// All GET requests will be handled with a redirect to the https HOST URL with
// a (301)Moved Permanently status code. All POST/PUT/DELETE requests will
// return an error specifying that a secure connection is required. This is done
// in order to avoid missing data when performing a redirect from
// POST -> GET redirects.
func ForceSSL(h httputils.HandlerFunc) httputils.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			ctx = r.Context()
			cfg = ctx.Value("config").(*config.Config)

			isNotHTTPS = r.Header.Get("X-Forwarded-Proto") != "https"
			isGet      = r.Method == "GET"
		)

		if isNotHTTPS {
			if isGet {
				http.Redirect(w, r, cfg.App.HostURL, http.StatusMovedPermanently)
			} else {
				http.Error(w, "Unsafe connection not allowed!", http.StatusBadRequest)
			}

			return nil
		}

		return h(w, r)
	}
}
