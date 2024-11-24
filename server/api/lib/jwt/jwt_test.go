package jwt_test

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	jwt2 "server/api/lib/jwt"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtHandler := jwt2.NewJWTHandler() // Создаем экземпляр структуры

	token, err := jwtHandler.GenerateAccessToken(1, "admin")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtHandler := jwt2.NewJWTHandler() // Создаем экземпляр структуры

	token, err := jwtHandler.GenerateRefreshToken(1)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestExtractJWTFromHeader(t *testing.T) {
	jwtHandler := jwt2.NewJWTHandler() // Создаем экземпляр структуры

	t.Run("Valid header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer testtoken")
		token, err := jwtHandler.ExtractJWTFromHeader(req)

		assert.NoError(t, err)
		assert.Equal(t, "testtoken", token)
	})

	t.Run("Missing header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		_, err := jwtHandler.ExtractJWTFromHeader(req)

		assert.Error(t, err)
		assert.Equal(t, "no authorization header", err.Error())
	})

	t.Run("Invalid header format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidHeader")
		_, err := jwtHandler.ExtractJWTFromHeader(req)

		assert.Error(t, err)
		assert.Equal(t, "invalid authorization header", err.Error())
	})
}

func TestValidateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtHandler := jwt2.NewJWTHandler() // Создаем экземпляр структуры

	t.Run("Valid token", func(t *testing.T) {
		token, _ := jwtHandler.GenerateAccessToken(1, "admin")
		claims, err := jwtHandler.ValidateJWT(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, int64(1), claims.UserId)
		assert.Equal(t, "admin", claims.Role)
	})

	t.Run("Expired token", func(t *testing.T) {
		expiredClaims := &jwt2.Claims{
			UserId: 1,
			Role:   "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredToken, _ := token.SignedString([]byte("testsecret"))

		_, err := jwtHandler.ValidateJWT(expiredToken)

		assert.Error(t, err)
		assert.Equal(t, "token expired", err.Error())
	})
}
