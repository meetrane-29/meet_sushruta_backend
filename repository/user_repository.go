package repository

import (
	"errors"
	"strings"

	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	EmailExists(email string) (bool, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id uuid.UUID) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	// Normalize email to lowercase
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	// Generate UUID if not provided
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Create user in database
	if err := config.DB.Create(user).Error; err != nil {
		// Check if it's a unique constraint violation (email already exists)
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("email already registered")
		}
		return nil, err
	}

	return user, nil
}

// EmailExists checks if an email is already registered
func (r *userRepository) EmailExists(email string) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	var count int64
	if err := config.DB.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	var user model.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(user *model.User) (*model.User, error) {
	if user == nil || user.ID == uuid.Nil {
		return nil, errors.New("invalid user")
	}

	if err := config.DB.Model(user).Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (r *userRepository) DeleteUser(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid user ID")
	}

	if err := config.DB.Delete(&model.User{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
