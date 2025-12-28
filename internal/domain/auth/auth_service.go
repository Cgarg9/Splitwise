package auth

import (
	"context"
	"errors"
	"splitwise-clone/internal/domain/user"
	"splitwise-clone/internal/logger"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
)

// Service defines the interface for authentication business logic
type Service interface {
	SignUp(ctx context.Context, params SignUpParams) (*user.User, error)
	// Login(ctx context.Context, email, password string) (string, error) // For future implementation
}

// SignUpParams contains parameters for user registration
type SignUpParams struct {
	FirstName   string
	LastName    string
	DateOfBirth *time.Time
	Email       string
	Password    string
	PhoneNumber *string
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new auth service instance
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// SignUp handles user registration
func (s *service) SignUp(ctx context.Context, params SignUpParams) (*user.User, error) {
	log := logger.FromContext(ctx)

	log.Debug().Str("email", params.Email).Msg("Starting user signup process")

	// Check if user already exists
	exists, err := s.repo.UserExistsByEmail(ctx, params.Email)
	if err != nil {
		log.Error().Err(err).Str("email", params.Email).Msg("Failed to check if user exists")
		return nil, err
	}
	if exists {
		log.Debug().Str("email", params.Email).Msg("User already exists")
		return nil, ErrUserAlreadyExists
	}

	// Hash the password using bcrypt
	log.Debug().Msg("Hashing password")
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, err
	}

	// Create user in database
	createParams := user.CreateUserParams{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		PasswordHash: hashedPassword,
		PhoneNumber:  params.PhoneNumber,
		DateOfBirth:  params.DateOfBirth,
	}

	log.Debug().Str("email", params.Email).Msg("Creating user in database")
	newUser, err := s.repo.CreateUser(ctx, createParams)
	if err != nil {
		log.Error().Err(err).Str("email", params.Email).Msg("Failed to create user in database")
		return nil, err
	}

	log.Info().
		Str("user_id", newUser.ID.String()).
		Str("email", newUser.Email).
		Msg("User signup completed successfully")

	return newUser, nil
}

// hashPassword hashes a plain text password using bcrypt with SHA-256
func hashPassword(password string) (string, error) {
	// Using bcrypt with cost 12 (recommended for production)
	// Note: bcrypt internally uses a secure algorithm
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword compares a plain text password with a hashed password
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
