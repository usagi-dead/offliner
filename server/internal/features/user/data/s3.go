package data

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"server/internal/features/user"
	"server/internal/s3Storage"
	"sync"
)

type UserS3Client interface {
	UploadAvatar(avatarSmall []byte, avatarLarge []byte, userId int64) (*string, error)
	DeleteAvatar(userId int64) error
}

type UserS3 struct {
	storage *s3Storage.S3Storage
}

func NewUserS3(storage *s3Storage.S3Storage) *UserS3 {
	return &UserS3{
		storage: storage,
	}
}

func (u *UserS3) UploadAvatar(avatarSmall []byte, avatarLarge []byte, userId int64) (*string, error) {
	folderPath := fmt.Sprintf("avatars/%d/", userId)

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
		return
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
		return
	}()

	wg.Wait()

	if errSmall != nil || errLarge != nil {
		return nil, user.ErrInternal
	}

	folderURL := fmt.Sprintf("https://%s.%s/%s", u.storage.Bucket, u.storage.Endpoint, folderPath)
	return &folderURL, nil
}

func (u *UserS3) DeleteAvatar(userId int64) error {
	folderPath := fmt.Sprintf("avatars/%d/", userId)

	objectsToDelete := []string{folderPath + "52x52.webp", folderPath + "512x512.webp"}

	var objectsId []types.ObjectIdentifier
	objectsId = append(objectsId, types.ObjectIdentifier{Key: aws.String(objectsToDelete[0])})
	objectsId = append(objectsId, types.ObjectIdentifier{Key: aws.String(objectsToDelete[1])})

	deleteInput := &s3.DeleteObjectsInput{
		Bucket: &u.storage.Bucket,
		Delete: &types.Delete{
			Objects: objectsId,
			Quiet:   aws.Bool(false),
		},
	}

	_, err := u.storage.Client.DeleteObjects(context.Background(), deleteInput)

	return err
}
