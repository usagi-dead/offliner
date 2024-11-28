package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"server/internal/config"
	u "server/internal/features/user"
	"strings"
	"time"
)

type JWTService interface {
	GenerateAccessToken(userID int64, role string) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ExtractJWTFromHeader(r *http.Request) (string, error)
	ValidateJWT(token string) (*Claims, error)
}

type JWTHandler struct {
	cfg *config.JWTConfig
}

func NewJWTHandler(cfg *config.JWTConfig) *JWTHandler {
	return &JWTHandler{
		cfg: cfg,
	}
}

type Claims struct {
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTHandler) GenerateAccessToken(UserID int64, Role string) (string, error) {
	claims := &Claims{
		UserId: UserID,
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (j *JWTHandler) GenerateRefreshToken(UserID int64) (string, error) {
	claims := &Claims{
		UserId: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.RefreshExpire)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (j *JWTHandler) ExtractJWTFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", u.ErrNoAccessToken
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", u.ErrInvalidToken
	}

	return authHeader[7:], nil
}

func (j *JWTHandler) ValidateJWT(JWTToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(JWTToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "is expired") {
			return nil, u.ErrExpiredToken
		}
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, u.ErrInvalidToken
	}

	return claims, nil
}
