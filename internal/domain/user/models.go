package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose password hash in JSON
	PhoneNumber  *string    `json:"phone_number,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// CreateUserParams contains parameters for creating a new user
type CreateUserParams struct {
	FirstName    string
	LastName     string
	DateOfBirth  *time.Time
	Email        string
	PasswordHash string
	PhoneNumber  *string
}
