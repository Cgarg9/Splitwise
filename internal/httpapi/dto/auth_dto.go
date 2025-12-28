package dto

import "time"

// SignUpRequest represents the request body for user registration
type SignUpRequest struct {
	FirstName   string     `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string     `json:"last_name" validate:"required,min=2,max=100"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Email       string     `json:"email" validate:"required,email"`
	Password    string     `json:"password" validate:"required,min=8,max=72"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164"`
}

// SignUpResponse represents the response after successful registration
type SignUpResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

