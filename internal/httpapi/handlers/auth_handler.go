package handlers

import (
	"encoding/json"
	"net/http"
	"splitwise-clone/internal/domain/auth"
	"splitwise-clone/internal/httpapi/dto"
	"splitwise-clone/internal/logger"

	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService auth.Service
	validate    *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// SignUp handles user registration requests
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Get logger from context (includes trace ID)
	log := logger.FromContext(r.Context())

	var req dto.SignUpRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode signup request")
		respondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Validation failed for signup request")
		respondWithError(w, http.StatusBadRequest, "Validation failed", parseValidationErrors(err))
		return
	}

	log.Info().Str("email", req.Email).Msg("Processing signup request")

	// Call service to create user
	signUpParams := auth.SignUpParams{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}

	user, err := h.authService.SignUp(r.Context(), signUpParams)
	if err != nil {
		if err == auth.ErrUserAlreadyExists {
			log.Warn().Str("email", req.Email).Msg("Signup failed: user already exists")
			respondWithError(w, http.StatusConflict, "User with this email already exists", nil)
			return
		}

		// Note: Unique constraint violations are already handled by ErrUserAlreadyExists check above

		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
		return
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Msg("User created successfully")

	// Build response
	response := dto.SignUpResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string, details map[string]interface{}) {
	respondWithJSON(w, code, dto.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Details: details,
	})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Note: We can't get context here for trace ID logging
		// At this point headers are already sent, so we can only log the error
		// Consider passing context to this function if detailed error tracking is needed
		_ = err // Acknowledge the error but can't do much at this point
	}
}

// parseValidationErrors converts validator errors to a map
func parseValidationErrors(err error) map[string]interface{} {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]interface{})
		for _, fieldError := range validationErrors {
			errors[fieldError.Field()] = fieldError.Tag()
		}
		return errors
	}
	return nil
}
