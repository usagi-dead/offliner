package services

import (
	"golang.org/x/crypto/bcrypt"
)

func (uuc *UserUseCase) SignUp(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err := uuc.repo.CreateUser(email, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}
