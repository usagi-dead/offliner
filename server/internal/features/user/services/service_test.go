package services_test

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"server/api/lib/jwt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"server/internal/features/user"
	"server/internal/features/user/services"
)

type MockUserData struct {
	mock.Mock
	emailConfirmationStatus map[string]bool
	emailConfirmationCodes  map[string]string
}

func NewMockUserData() *MockUserData {
	return &MockUserData{
		emailConfirmationStatus: make(map[string]bool),
		emailConfirmationCodes:  make(map[string]string),
	}
}

func (m *MockUserData) CongirmEmail(email string) error {
	if _, exists := m.emailConfirmationCodes[email]; !exists {
		return fmt.Errorf("no confirmation code found for email: %s", email)
	}
	m.emailConfirmationStatus[email] = true
	return nil
}

func (m *MockUserData) IsEmailConfirmed(email string) (bool, error) {
	confirmed, exists := m.emailConfirmationStatus[email]
	if !exists {
		return false, user.ErrUserNotFound
	}
	return confirmed, nil
}

func (m *MockUserData) SaveEmailConfirmedCode(email string, code string) error {
	m.emailConfirmationStatus[email] = false
	m.emailConfirmationCodes[email] = code
	return nil
}

func (m *MockUserData) GetEmailConfirmedCode(email string) (string, error) {

	code, exists := m.emailConfirmationCodes[email]
	if !exists {
		return "", fmt.Errorf("no confirmation code found for email: %s", email)
	}
	return code, nil
}

func (m *MockUserData) CreateUser(email string, hashedPassword string) (int64, error) {
	args := m.Called(email, hashedPassword)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserData) CreateOauthUser(user *user.User) (int64, error) {
	args := m.Called(user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserData) GetUserByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	// Возвращаем указатель на объект пользователя, даже если он пустой
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserData) SaveStateCode(state string) error {
	return nil
}

func (m *MockUserData) VerifyStateCode(state string) (bool, error) {
	return false, nil
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID int64, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateJWT(token string) (*jwt.Claims, error) {
	return nil, nil
}

func (m *MockJWTService) ExtractJWTFromHeader(r *http.Request) (string, error) {
	return "", nil
}

type MockEmailSender struct{}

func (m *MockEmailSender) SendConfirmEmail(email string, code string) error {
	// Логика мока, просто возвращаем успех
	return nil
}

func TestUserUseCase_SignUp(t *testing.T) {
	mockRepo := new(MockUserData)
	mockJWT := new(MockJWTService)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	es := new(MockEmailSender)
	userUseCase := services.NewUserUseCase(mockRepo, mockJWT, logger, es)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		password := "securepassword"
		userID := int64(1)

		// Настройка мока репозитория
		mockRepo.On("CreateUser", email, mock.Anything).Return(userID, nil)

		// Выполнение метода
		err := userUseCase.SignUp(email, password)

		// Проверки
		assert.NoError(t, err)

		mockRepo.AssertCalled(t, "CreateUser", email, mock.Anything)
	})
}

func TestUserUseCase_SignIn(t *testing.T) {
	mockRepo := new(MockUserData)
	mockJWT := new(MockJWTService)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	es := new(MockEmailSender)
	userUseCase := services.NewUserUseCase(mockRepo, mockJWT, logger, es)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		password := "securepassword"
		hashedPasswordbyte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hashedPassword := string(hashedPasswordbyte)
		userID := int64(1)
		role := "user"
		accessToken := "access-token"
		refreshToken := "refresh-token"

		// Настройка мока репозитория
		mockRepo.On("GetUserByEmail", email).Return(&user.User{
			UserId:         userID,
			HashedPassword: &hashedPassword,
			Role:           role,
			VerifiedEmail:  true,
		}, nil)

		// Настройка мока JWT
		mockJWT.On("GenerateAccessToken", userID, role).Return(accessToken, nil)
		mockJWT.On("GenerateRefreshToken", userID).Return(refreshToken, nil)

		// Выполнение метода
		act, ref, err := userUseCase.SignIn(email, password)

		// Проверки
		assert.NoError(t, err)
		assert.Equal(t, accessToken, act)
		assert.Equal(t, refreshToken, ref)

		mockRepo.AssertCalled(t, "GetUserByEmail", email)
		mockJWT.AssertCalled(t, "GenerateAccessToken", userID, role)
		mockJWT.AssertCalled(t, "GenerateRefreshToken", userID)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		email := "test@example.com"
		password := "wrongpassword"

		// Настройка мока репозитория
		mockRepo.On("GetUserByEmail", email).Return(&user.User{
			HashedPassword: nil,
		}, nil)

		// Выполнение метода
		act, ref, err := userUseCase.SignIn(email, password)

		// Проверки
		assert.Empty(t, act)
		assert.Empty(t, ref)
		assert.Error(t, err)
		assert.Equal(t, err, user.ErrUserNotFound)

		mockRepo.AssertCalled(t, "GetUserByEmail", email)
	})

	t.Run("email not verified", func(t *testing.T) {
		email := "test@example.com"
		password := "securepassword"
		hashedPasswordbyte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hashedPassword := string(hashedPasswordbyte)

		mockRepo.On("GetUserByEmail", email).Return(&user.User{
			UserId:         int64(2),
			HashedPassword: &hashedPassword,
			Role:           "user",
			VerifiedEmail:  false,
		}, nil)

		act, ref, err := userUseCase.SignIn(email, password)

		// Проверки
		assert.Empty(t, act)
		assert.Empty(t, ref)
		assert.Error(t, err)
		assert.Equal(t, user.ErrEmailNotConfirmed, err)

		mockRepo.AssertCalled(t, "GetUserByEmail", email)
	})
}

func TestUserUseCase_EmailConfirmed(t *testing.T) {
	mockRepo := NewMockUserData()
	mockJWT := new(MockJWTService)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	es := new(MockEmailSender)
	userUseCase := services.NewUserUseCase(mockRepo, mockJWT, logger, es)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		code := "asdfKD"
		_ = mockRepo.SaveEmailConfirmedCode(email, code)

		err := userUseCase.EmailConfirmed(email, code)

		assert.NoError(t, err)
	})

	t.Run("email already confirmed", func(t *testing.T) {
		email := "test@example.com"
		code := "asdfKD"
		_ = mockRepo.SaveEmailConfirmedCode(email, code)
		_ = mockRepo.CongirmEmail(email)

		err := userUseCase.EmailConfirmed(email, code)

		assert.Error(t, err)
		assert.Equal(t, user.ErrEmailAlreadyConfirmed, err)
	})

	t.Run("user not found", func(t *testing.T) {
		email := "test@example.com"
		code := "asdfKD"

		err := userUseCase.EmailConfirmed(email, code)

		assert.Error(t, err)
		assert.Equal(t, user.ErrUserNotFound, err)
	})

	t.Run("invalid confirm code", func(t *testing.T) {
		email := "test@example.com"
		code := "asdfKD"

		_ = mockRepo.SaveEmailConfirmedCode(email, code)

		err := userUseCase.EmailConfirmed(email, "asdfasqewr")

		assert.Error(t, err)
		assert.Equal(t, user.ErrInvalidConfirmCode, err)
	})
}

func TestUserUseCase_SendEmailForConfirmed(t *testing.T) {
	mockRepo := NewMockUserData()
	mockJWT := new(MockJWTService)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	es := new(MockEmailSender)
	userUseCase := services.NewUserUseCase(mockRepo, mockJWT, logger, es)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"

		mockRepo.emailConfirmationStatus[email] = false

		err := userUseCase.SendEmailForConfirmed(email)

		code, err := mockRepo.GetEmailConfirmedCode(email)

		assert.NoError(t, err)
		assert.NotEmpty(t, code)
	})

	t.Run("email already confirmed", func(t *testing.T) {
		email := "test@example.com"
		code := "asdfKD"

		_ = mockRepo.SaveEmailConfirmedCode(email, code)
		_ = mockRepo.CongirmEmail(email)

		err := userUseCase.SendEmailForConfirmed(email)

		assert.Error(t, err)
		assert.Equal(t, user.ErrEmailAlreadyConfirmed, err)
	})
}

func TestUserUseCase_GenerateEmailCode(t *testing.T) {
	code, err := services.GenerateEmailCode()
	if err != nil {
		t.Fatalf("GenerateEmailCode() returned an error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("Expected code length to be 6, got %d", len(code))
	}

	const CharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range code {
		if !isCharInSet(char, CharSet) {
			t.Errorf("Code contains invalid character: %v", char)
		}
	}

	testRandomness(t, 1000)
}

func isCharInSet(char rune, set string) bool {
	for _, s := range set {
		if char == s {
			return true
		}
	}
	return false
}

func testRandomness(t *testing.T, iterations int) {
	seen := make(map[string]bool)
	for i := 0; i < iterations; i++ {
		code, err := services.GenerateEmailCode()
		if err != nil {
			t.Fatalf("GenerateEmailCode() returned an error during randomness test: %v", err)
		}
		if seen[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		seen[code] = true
	}
}
