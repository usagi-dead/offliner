package services

import (
	"crypto/rand"
	"errors"
	"math/big"
	u "server/internal/features/user"
)

func (uuc *UserUseCase) SendEmailForConfirmed(email string) error {
	ok, err := uuc.repo.IsEmailConfirmed(email)
	if err != nil {
		if errors.Is(err, u.ErrUserNotFound) {
			return u.ErrUserNotFound
		}
		return err
	}
	if ok {
		return u.ErrEmailAlreadyConfirmed
	}

	code, err := GenerateEmailCode()
	if err != nil {
		return err
	}

	if err := uuc.repo.SaveEmailConfirmedCode(email, code); err != nil {
		return err
	}

	if err := uuc.ess.SendConfirmEmail(code, email); err != nil {
		return err
	}

	return nil
}

func GenerateEmailCode() (string, error) {
	const CharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(CharSet))))
		if err != nil {
			return "", err
		}
		code[i] = CharSet[index.Int64()]
	}
	return string(code), nil
}
