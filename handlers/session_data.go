package handlers

import "time"

// SessionData describes the session cookie for all users.
type SessionData struct {
	UserID    int
	Email     string
	ExpiresAt time.Time
}

// IsInvalid checks wether the data is in the correct state.
func (s *SessionData) IsInvalid() bool {
	if s.UserID == 0 {
		return true
	}

	if s.ExpiresAt.IsZero() {
		return true
	}

	return false
}
