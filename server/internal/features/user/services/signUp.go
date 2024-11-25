package services

import (
	"golang.org/x/crypto/bcrypt"
)

func (uuc *UserUseCase) SignUp(email string, password string) (string, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	userID, err := uuc.repo.CreateUser(email, string(hashedPassword))
	if err != nil {
		return "", "", err
	}

	accessToken, err := uuc.jwt.GenerateAccessToken(userID, "user")
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uuc.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
