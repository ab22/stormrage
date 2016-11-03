package services

import (
	"errors"
	"fmt"
)

var (
	// ErrRecordNotFound indicates that a query returned no rows.
	ErrRecordNotFound = errors.New("record not found")
)

// ErrUserAlreadyExists contains information about the user that already
// exists in the database.
type ErrUserAlreadyExists string

func (e ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("could not create user: user [%v] already exists in the database!", string(e))
}

// ErrExpiredToken indicates that the token already expired.
type ErrExpiredToken struct{}

func (e *ErrExpiredToken) Error() string {
	return fmt.Sprintf("token expired")
}
