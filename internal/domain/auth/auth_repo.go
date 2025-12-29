package auth

import (
	"context"
	"splitwise-clone/internal/domain/user"
	"splitwise-clone/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for authentication-related database operations
type Repository interface {
	CreateUser(ctx context.Context, params user.CreateUserParams) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new auth repository instance
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

// CreateUser inserts a new user into the database
func (r *repository) CreateUser(ctx context.Context, params user.CreateUserParams) (*user.User, error) {
	log := logger.FromContext(ctx)

	query := `
		INSERT INTO users (first_name, last_name, date_of_birth, email, password_hash, phone_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, last_name, date_of_birth, email, password_hash, phone_number, created_at, updated_at, deleted_at
	`

	log.Debug().Str("email", params.Email).Msg("Executing CreateUser query")

	var u user.User
	err := r.db.QueryRow(ctx, query,
		params.FirstName,
		params.LastName,
		params.DateOfBirth,
		params.Email,
		params.PasswordHash,
		params.PhoneNumber,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.DateOfBirth,
		&u.Email,
		&u.PasswordHash,
		&u.PhoneNumber,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if err != nil {
		log.Error().Err(err).Str("email", params.Email).Msg("Failed to insert user into database")
		return nil, err
	}

	log.Debug().Str("user_id", u.ID.String()).Msg("User inserted successfully")
	return &u, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	log := logger.FromContext(ctx)

	query := `
		SELECT id, first_name, last_name, date_of_birth, email, password_hash, phone_number, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	log.Debug().Str("email", email).Msg("Executing GetUserByEmail query")

	var u user.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.DateOfBirth,
		&u.Email,
		&u.PasswordHash,
		&u.PhoneNumber,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if err != nil {
		log.Debug().Err(err).Str("email", email).Msg("User not found or query failed")
		return nil, err
	}

	log.Debug().Str("user_id", u.ID.String()).Msg("User retrieved successfully")
	return &u, nil
}

// UserExistsByEmail checks if a user with the given email already exists
func (r *repository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	log := logger.FromContext(ctx)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

	log.Debug().Str("email", email).Msg("Checking if user exists")

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to check user existence")
		return false, err
	}

	log.Debug().Str("email", email).Bool("exists", exists).Msg("User existence check completed")
	return exists, nil
}
