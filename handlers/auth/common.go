package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/handlers"
	"github.com/ab22/stormrage/handlers/httputils"
	"github.com/gorilla/sessions"
)

// CheckAuth asumes that the ValidateAuth decorator called this function
// because the session was validated successfully.
func (h *handler) CheckAuth(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Login does basic email/password login.
// Checks:
// 		- User must exist
//		- Passwords match
//		- User's status is Active
// If the checks pass, it sets up a session cookie.
func (h *handler) Login(w http.ResponseWriter, r *http.Request) error {
	var (
		err         error
		ctx         = r.Context()
		cookieStore = ctx.Value("cookieStore").(*sessions.CookieStore)
		cfg         = ctx.Value("config").(*config.Config)

		loginForm struct {
			Username string
			Password string
		}
	)

	if err = httputils.DecodeJSON(r.Body, &loginForm); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "")
		return nil
	}

	user, err := h.authService.BasicAuth(loginForm.Username, loginForm.Password)

	if err != nil {
		return err
	} else if user == nil {
		var errorMsg = fmt.Sprintf(
			"Failed login attempt with user [%s] from IP [%s]",
			loginForm.Username,
			r.RemoteAddr,
		)
		log.Println(errorMsg)

		httputils.WriteError(w, http.StatusUnauthorized, "Usuario/Clave inválidos!")
		return nil
	}

	session, err := cookieStore.New(r, cfg.SessionCookieName)
	session.Values["data"] = &handlers.SessionData{
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(cfg.SessionLifeTime),
	}

	return session.Save(r, w)
}

// Logout does basic session logout.
func (h *handler) Logout(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx          = r.Context()
		cfg          = ctx.Value("config").(*config.Config)
		cookieStore  = ctx.Value("cookieStore").(*sessions.CookieStore)
		session, err = cookieStore.Get(r, cfg.SessionCookieName)
	)

	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}
