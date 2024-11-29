package data

import (
	u "server/internal/features/user"
)

type UserQuery struct {
	db UserDBClient
	ch UserCacheClient
	s3 UserS3Client
}

func NewUserQuery(db UserDBClient, ch UserCacheClient, s3 UserS3Client) u.UserData {
	return &UserQuery{
		db: db,
		ch: ch,
		s3: s3,
	}
}

func (uq *UserQuery) CreateUser(email string, hashedPassword string) (int64, error) {
	return uq.db.CreateUser(email, hashedPassword)
}

func (uq *UserQuery) CreateOauthUser(user *u.User) (int64, error) {
	return uq.db.CreateOauthUser(user)
}

func (uq *UserQuery) GetUserByEmail(email string) (*u.User, error) {
	return uq.db.GetUserByEmail(email)
}

func (uq *UserQuery) GetUserById(userId int64) (*u.User, error) {
	return uq.db.GetUserById(userId)
}

func (uq *UserQuery) SaveStateCode(state string) error {
	return uq.ch.SaveStateCode(state)
}

func (uq *UserQuery) VerifyStateCode(state string) (bool, error) {
	return uq.ch.VerifyStateCode(state)
}

func (uq *UserQuery) ConfirmEmail(email string) error {
	return uq.db.ConfirmEmail(email)
}

func (uq *UserQuery) IsEmailConfirmed(email string) (bool, error) {
	return uq.db.IsEmailConfirmed(email)
}

func (uq *UserQuery) SaveEmailConfirmedCode(email string, code string) error {
	return uq.ch.SaveEmailConfirmedCode(email, code)
}

func (uq *UserQuery) GetEmailConfirmedCode(email string) (string, error) {
	return uq.ch.GetEmailConfirmedCode(email)
}

func (uq *UserQuery) UploadAvatar(avatarSmall []byte, avatarLarge []byte, userId string) (string, error) {
	return uq.s3.UploadAvatar(avatarSmall, avatarLarge, userId)
}
