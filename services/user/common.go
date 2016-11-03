package user

import (
	"github.com/ab22/stormrage/models"
	"github.com/ab22/stormrage/services"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Searches for a User by Email.
// Returns *models.User instance if it finds it, or nil otherwise.
func (s *service) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := s.db.
		Where("email = ?", email).
		First(user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return user, nil
}

// Searches for a User by Username.
// Returns *models.User instance if it finds it, or nil otherwise.
func (s *service) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}

	err := s.db.
		Where("username = ?", username).
		First(user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return user, nil
}

// Encrypts a password with the default password hasher (bcrypt).
// Returns the hashed password []byte and an error.
func (s *service) EncryptPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// Compares if the hashed password equals the plain text password.
func (s *service) ComparePasswords(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

// Checks if a user with that email already exists in the database. If it does,
// it returns an error, else it hashes the password, saves the new user
// and returns the user.
func (s *service) CreateUser(email, password, firstName, lastName string, status Status) (*models.User, error) {
	var err error

	result, err := s.FindByEmail(email)
	if err != nil {
		return nil, err
	} else if result != nil {
		return nil, services.ErrUserAlreadyExists(email)
	}

	hashedPassword, err := s.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Status:    int(status),
	}

	err = s.db.Create(&user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ActivateUser searches for a user in the database by email and updates
// it's status to Active.
// Note: To avoid issues with possible future banned users, the ActivateUser
// makes sure that that user selected is in state Unconfirmed (state that is
// set up only when creating a user).
func (s *service) ActivateUser(email string) error {
	result := s.db.
		Table("users").
		Where("email = ?", email).
		Where("status = ?", int(Unconfirmed)).
		Update("status", int(Active))

	if err := result.Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else if result.RowsAffected == 0 {
		return services.ErrRecordNotFound
	}

	return nil
}

// ChangePassword finds a user in the database by email and changes it's
// password.
func (s *service) ChangePassword(username, password string) error {
	hashedPassword, err := s.EncryptPassword(password)

	result := s.db.
		Table("users").
		Where("username = ?", username).
		Update("password", string(hashedPassword))

	if err = result.Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else if result.RowsAffected == 0 {
		return services.ErrRecordNotFound
	}

	return nil
}
