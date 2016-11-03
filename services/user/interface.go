package user

import (
	"github.com/ab22/stormrage/models"
	"github.com/jinzhu/gorm"
)

// Service interface describes all functions that must be implemented.
type Service interface {
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	EncryptPassword(password string) ([]byte, error)
	ComparePasswords(hashedPassword []byte, password string) bool
	CreateUser(email, password, firstName, lastName string, status Status) (*models.User, error)
	ActivateUser(email string) error
	ChangePassword(email, password string) error
}

// Status defines statuses for the User model.
type Status int

// Defines all user statuses.
const (
	Unconfirmed Status = iota
	Active
)

// Contains all of the logic for the User model.
type service struct {
	db *gorm.DB
}

// NewService initialization.
func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
