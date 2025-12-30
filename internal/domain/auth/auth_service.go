package auth

import (
	"context"
	"errors"
	"os"
	"splitwise-clone/internal/domain/user"
	"splitwise-clone/internal/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrTokenGeneration   = errors.New("failed to generate token")
)

// JWT configuration
var (
	jwtSecret = []byte(getJWTSecret())
	jwtExpiry = 24 * time.Hour // Token expires in 24 hours
)

// Claims represents the JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// Service defines the interface for authentication business logic
type Service interface {
	SignUp(ctx context.Context, params SignUpParams) (*user.User, error)
	Login(ctx context.Context, email, password string) (string, error)
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

// Login handles user authentication
func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	log := logger.FromContext(ctx)

	log.Debug().Str("email", email).Msg("Starting user login process")

	// Check if user exists
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
		return "", err
	}
	if user == nil {
		log.Debug().Str("email", email).Msg("User not found")
		return "", ErrUserNotFound
	}

	// Verify Password
	err = verifyPassword(user.PasswordHash, password)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to verify password")
		return "", ErrInvalidPassword
	}

	log.Info().Str("email", email).Msg("User authenticated successfully")

	// Generate JWT Token
	token, err := generateJWTToken(user.ID)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to generate JWT token")
		return "", err
	}

	return token, nil
}

// generateJWTToken generates a JWT token for the given user ID
func generateJWTToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(jwtExpiry)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "splitwise-clone",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return tokenString, nil
}

// getJWTSecret retrieves JWT secret from environment or uses default for development
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret for development only
		// In production, JWT_SECRET environment variable MUST be set
		return "dev-secret-key-change-in-production"
	}
	return secret
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
// This will be used in the Login implementation
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
