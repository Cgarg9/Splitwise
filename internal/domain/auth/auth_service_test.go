package auth

import (
	"context"
	"splitwise-clone/internal/domain/user"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, params user.CreateUserParams) (*user.User, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func createTestSignUpParams() SignUpParams {
	return SignUpParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "securePassword123",
	}
}

// createTestUser creates a test user entity
func createTestUser() *user.User {
	now := time.Now()
	return &user.User{
		ID:           uuid.New(),
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "$2a$12$hashedpassword",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func TestSignUp_Success(t *testing.T) {

	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	params := createTestSignUpParams()
	expectedUser := createTestUser()

	// define what mocks should return
	mockRepo.On("UserExistsByEmail", ctx, params.Email).Return(false, nil)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("user.CreateUserParams")).Return(expectedUser, nil)

	result, err := service.SignUp(ctx, params)

	// assertions - verify the result
	assert.NoError(t, err, "SignUp should not return an error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, expectedUser.ID, result.ID, "user IDs should match")
	assert.Equal(t, expectedUser.Email, result.Email, "user emails should match")
	assert.Equal(t, expectedUser.FirstName, result.FirstName, "user first names should match")
	assert.Equal(t, expectedUser.LastName, result.LastName, "user last names should match")

	mockRepo.AssertExpectations(t)
}
