package services

import (
	"errors"
	u "server/internal/features/user"
)

func (uus *UserUseCase) EmailConfirmed(email string, code string) error {
	ok, err := uus.repo.IsEmailConfirmed(email)
	if err != nil {
		if errors.Is(err, u.ErrUserNotFound) {
			return u.ErrUserNotFound
		}
		return err
	}
	if ok {
		return u.ErrEmailAlreadyConfirmed
	}

	realCode, err := uus.repo.GetEmailConfirmedCode(email)
	if err != nil {
		return err
	}

	if realCode != code {
		return u.ErrInvalidConfirmCode
	}

	if err := uus.repo.CongirmEmail(email); err != nil {
		return err
	}

	return nil
}
