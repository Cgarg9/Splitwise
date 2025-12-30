package auth

import (
	"context"
	"errors"
	"splitwise-clone/internal/domain/user"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestSignUp_OptionalFields(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	phone := "+1234567890"
	dob := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
	params := SignUpParams{
		FirstName:   "Jane",
		LastName:    "Smith",
		Email:       "jane.smith@example.com",
		Password:    "securePassword123",
		PhoneNumber: &phone,
		DateOfBirth: &dob,
	}

	expectedUser := createTestUser()
	expectedUser.DateOfBirth = &dob
	expectedUser.PhoneNumber = &phone
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
	assert.Equal(t, expectedUser.DateOfBirth, result.DateOfBirth, "user date of births should match")
	assert.Equal(t, expectedUser.PhoneNumber, result.PhoneNumber, "user phone numbers should match")

	mockRepo.AssertExpectations(t)
}

func TestSignUp_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	params := createTestSignUpParams()

	// define what mocks should return
	mockRepo.On("UserExistsByEmail", ctx, params.Email).Return(true, nil)

	result, err := service.SignUp(ctx, params)

	// assertions - verify the result
	assert.Error(t, err, "SignUp should return an error")
	assert.Nil(t, result, "Result should be nil")
	assert.Equal(t, ErrUserAlreadyExists, err, "Error should be ErrUserAlreadyExists")

	mockRepo.AssertNotCalled(t, "CreateUser")
	mockRepo.AssertExpectations(t)
}

func TestSignUp_UserExistsCheckError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	params := createTestSignUpParams()

	// define what mocks should return
	dbError := errors.New("database error")
	mockRepo.On("UserExistsByEmail", ctx, params.Email).Return(false, dbError)

	result, err := service.SignUp(ctx, params)

	// assertions - verify the result
	assert.Error(t, err, "SignUp should return an error")
	assert.Nil(t, result, "Result should be nil")
	assert.Equal(t, dbError, err, "Error should match the database error")

	mockRepo.AssertNotCalled(t, "CreateUser")
	mockRepo.AssertExpectations(t)
}

func TestSignUp_CreateUserError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	params := createTestSignUpParams()

	dbError := errors.New("Database insert error")

	mockRepo.On("UserExistsByEmail", ctx, params.Email).Return(false, nil)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("user.CreateUserParams")).Return(nil, dbError)

	result, err := service.SignUp(ctx, params)

	assert.Error(t, err, "SignUp should return an error")
	assert.Nil(t, result, "Result should be nil")
	assert.Equal(t, err, dbError, "Error should match the database error")

	mockRepo.AssertExpectations(t)

}

func TestSignUp_PasswordIsHashed(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	params := createTestSignUpParams()
	expectedUser := createTestUser()

	var capturedParams user.CreateUserParams

	mockRepo.On("UserExistsByEmail", ctx, params.Email).Return(false, nil)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("user.CreateUserParams")).
		Run(func(args mock.Arguments) {
			capturedParams = args.Get(1).(user.CreateUserParams)
		}).
		Return(expectedUser, nil)

	result, err := service.SignUp(ctx, params)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEqual(t, params.Password, capturedParams.PasswordHash, "Password should be hashed")
	assert.NotEmpty(t, capturedParams.PasswordHash, "Hash cannot be empty")

	err = bcrypt.CompareHashAndPassword(
		[]byte(capturedParams.PasswordHash),
		[]byte(params.Password),
	)
	assert.NoError(t, err, "Password hash should match original password")

	mockRepo.AssertExpectations(t)

}

func TestSignUp_MultipleScenarios(t *testing.T) {
	// Table-driven test for multiple scenarios
	tests := []struct {
		name          string
		params        SignUpParams
		userExists    bool
		existsError   error
		createError   error
		expectedError error
		shouldSucceed bool
	}{
		{
			name:          "successful signup",
			params:        createTestSignUpParams(),
			userExists:    false,
			existsError:   nil,
			createError:   nil,
			expectedError: nil,
			shouldSucceed: true,
		},
		{
			name:          "user already exists",
			params:        createTestSignUpParams(),
			userExists:    true,
			existsError:   nil,
			createError:   nil,
			expectedError: ErrUserAlreadyExists,
			shouldSucceed: false,
		},
		{
			name:          "database error on exists check",
			params:        createTestSignUpParams(),
			userExists:    false,
			existsError:   errors.New("db error"),
			createError:   nil,
			expectedError: errors.New("db error"),
			shouldSucceed: false,
		},
		{
			name:          "database error on create",
			params:        createTestSignUpParams(),
			userExists:    false,
			existsError:   nil,
			createError:   errors.New("insert error"),
			expectedError: errors.New("insert error"),
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			service := NewService(mockRepo)
			ctx := context.Background()

			mockRepo.On("UserExistsByEmail", ctx, tt.params.Email).
				Return(tt.userExists, tt.existsError)

			if !tt.userExists && tt.existsError == nil {
				if tt.createError != nil {
					mockRepo.On("CreateUser", ctx, mock.AnythingOfType("user.CreateUserParams")).
						Return(nil, tt.createError)
				} else {
					mockRepo.On("CreateUser", ctx, mock.AnythingOfType("user.CreateUserParams")).
						Return(createTestUser(), nil)
				}
			}

			result, err := service.SignUp(ctx, tt.params)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
