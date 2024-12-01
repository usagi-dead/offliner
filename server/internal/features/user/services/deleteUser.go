package services

func (uuc *UserUseCase) DeleteUser(userId int64) error {
	return uuc.repo.DeleteUser(userId)
}
