package dto

import "time"

// SignUpRequest represents the request body for user registration
type SignUpRequest struct {
	FirstName   string     `json:"first_name" validate:"required,min=2,max=100" example:"John"`
	LastName    string     `json:"last_name" validate:"required,min=2,max=100" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01T00:00:00Z"`
	Email       string     `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password    string     `json:"password" validate:"required,min=8,max=72" example:"P@ssw0rd"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+1234567890"`
}

// SignUpResponse represents the response after successful registration
type SignUpResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error" example:"Bad Request"`
	Message string                 `json:"message" example:"Invalid request body"`
	Details map[string]interface{} `json:"details,omitempty"`
}
