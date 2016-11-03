package auth

import (
	"github.com/ab22/stormrage/models"
	"github.com/ab22/stormrage/services/user"
	"github.com/jinzhu/gorm"
)

// Service interface describes all functions that must be implemented.
type Service interface {
	BasicAuth(email, password string) (*models.User, error)
}

// service contains all of the logic for the systems authentications.
type service struct {
	db          *gorm.DB
	userService user.Service
}

// NewService initialization.
func NewService(db *gorm.DB, userService user.Service) Service {
	return &service{
		db:          db,
		userService: userService,
	}
}
