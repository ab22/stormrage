package auth

import "net/http"

// CheckAuth asumes that the ValidateAuth decorator called this function
// because the session was validated successfully.
func (h *handler) CheckAuth(w http.ResponseWriter, r *http.Request) error {
	return nil
}
