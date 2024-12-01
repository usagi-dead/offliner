package services

import u "server/internal/features/user"

func (uuc *UserUseCase) GetUser(userId int64) (*u.User, error) {
	user, err := uuc.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
