package services_test

import (
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"server/api/lib/emailsender"
	"server/api/lib/jwt"
	"server/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"server/internal/features/user"
	"server/internal/features/user/services"
)

type MockUserData struct {
	mock.Mock
}

func (m *MockUserData) CongirmEmail(email string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserData) IsEmailConfirmed(email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserData) SaveEmailConfirmedCode(email string, code string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserData) GetEmailConfirmedCode(email string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserData) CreateUser(email string, hashedPassword string) (int64, error) {
	args := m.Called(email, hashedPassword)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserData) CreateOauthUser(user *user.User) (int64, error) {
	return 0, nil
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

func TestUserUseCase_SignUp(t *testing.T) {
	mockRepo := new(MockUserData)
	mockJWT := new(MockJWTService)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	es, _ := emailsender.New(config.SMTPConfig{
		Host:     "smtp.yandex.ru",
		Port:     465,
		Username: "OfflinerMen@yandex.by",
	})
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
	es, _ := emailsender.New(config.SMTPConfig{
		Host:     "smtp.yandex.ru",
		Port:     465,
		Username: "OfflinerMen@yandex.by",
	})
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

func TestUserUseCase_oAuth(t *testing.T) {

}

func TestUserUseCase_EmailConfirmed(t *testing.T) {

}

func TestUserUseCase_SendEmailForConfirmed(t *testing.T) {

}

func TestUserUseCase_GenerateEmailCode(t *testing.T) {

}
