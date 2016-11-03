package auth

import (
	"github.com/ab22/stormrage/models"
	"github.com/ab22/stormrage/services/user"
)

// Basic username/password authentication. BasicAuth checks if the user exists,
// checks if the passwords match and if the user's state is active.
func (s *service) BasicAuth(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, nil
	}

	u, err := s.userService.FindByUsername(username)

	if err != nil {
		return nil, err
	} else if u == nil || u.Status != int(user.Active) {
		return nil, nil
	}

	match := s.userService.ComparePasswords([]byte(u.Password), password)
	if !match {
		return nil, nil
	}

	return u, nil
}
