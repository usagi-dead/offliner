package services

import (
	"net/http"
	u "server/internal/features/user"
)

func (uuc *UserUseCase) RefreshToken(r *http.Request) (string, error) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		return "", u.ErrNoRefreshToken
	}

	claims, err := uuc.jwt.ValidateJWT(refreshToken.Value)
	if err != nil {
		return "", err
	}

	user, err := uuc.repo.GetUserById(claims.UserId)
	if err != nil {
		return "", err
	}

	accessToken, err := uuc.jwt.GenerateAccessToken(claims.UserId, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
