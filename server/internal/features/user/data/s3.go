package data

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"server/internal/features/user"
	"server/internal/s3Storage"
	"sync"
)

type UserS3Client interface {
	UploadAvatar(avatarSmall []byte, avatarLarge []byte, userId string) (string, error)
	DownloadAvatar(userId string) (string, error)
	DeleteAvatar(userId string) error
}

type UserS3 struct {
	storage *s3Storage.S3Storage
}

func NewUserS3(storage *s3Storage.S3Storage) *UserS3 {
	return &UserS3{
		storage: storage,
	}
}

func (u *UserS3) UploadAvatar(avatarSmall []byte, avatarLarge []byte, userId string) (string, error) {
	folderPath := fmt.Sprintf("avatars/%s/", userId)

	objectKeySmall := folderPath + "52x52.webp"
	objectKeyLarge := folderPath + "512x512.webp"

	var wg sync.WaitGroup
	var errSmall, errLarge error

	wg.Add(1)
	go func() {
		defer wg.Done()

		uploadInputSmall := &s3.PutObjectInput{
			Bucket:      &u.storage.Bucket,
			Key:         &objectKeySmall,
			Body:        bytes.NewReader(avatarSmall),
			ContentType: aws.String("image/webp"),
		}

		_, errSmall = u.storage.Client.PutObject(context.TODO(), uploadInputSmall)
		if errSmall != nil {
			log.Printf("unable to upload small avatar to S3, %v", errSmall)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		uploadInputLarge := &s3.PutObjectInput{
			Bucket:      &u.storage.Bucket,
			Key:         &objectKeyLarge,
			Body:        bytes.NewReader(avatarLarge),
			ContentType: aws.String("image/webp"),
		}

		_, errLarge = u.storage.Client.PutObject(context.TODO(), uploadInputLarge)
		if errLarge != nil {
			log.Printf("unable to upload large avatar to S3, %v", errLarge)
			return
		}
	}()

	wg.Wait()

	if errSmall != nil || errLarge != nil {
		return "", user.ErrInternal
	}

	folderURL := fmt.Sprintf("https://%s.%s/%s", u.storage.Bucket, u.storage.Endpoint, folderPath)
	return folderURL, nil
}

func (u *UserS3) DownloadAvatar(userId string) (string, error) {
	//TODO: implement me
	return "", nil
}

func (u *UserS3) DeleteAvatar(userId string) error {
	//TODO: implement me
	return nil
}
