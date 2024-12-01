package services

import (
	"log"
	"mime/multipart"
	avatarManager "server/api/lib/avatarMenager"
	u "server/internal/features/user"
)

func (uuc *UserUseCase) UpdateUser(avatar *multipart.File, user *u.User, resetAvatar bool) error {

	if resetAvatar {
		if avatar == nil {
			da := "https://useravatars.storage-173.s3hoster.by/default"
			if err := uuc.repo.DeleteAvatar(user.UserId); err != nil {
				log.Print(err)
			}
			user.AvatarUrl = &da
		} else {
			smallAvatar, largeAvatar, err := avatarManager.ParsingAvatarImage(avatar)
			if err != nil {
				return err
			}

			user.AvatarUrl, err = uuc.repo.UploadAvatar(smallAvatar, largeAvatar, user.UserId)
			if err != nil {
				return err
			}
		}
	}

	if err := uuc.repo.UpdateUser(user); err != nil {
		return u.ErrInternal
	}

	return nil
}
