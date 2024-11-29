package s3Storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	config2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"server/internal/config"
	"time"
)

type S3Storage struct {
	Client   *s3.Client
	Bucket   string
	Endpoint string
}

func NewS3Storage(config config.S3Config) (*S3Storage, error) {
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		return nil, errors.New("S3 env not set")
	}

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL: "https://" + config.Endpoint,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config2.LoadDefaultConfig(context.Background(),
		config2.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config2.WithRegion(config.Region),
		config2.WithEndpointResolver(customResolver),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	// Проверяем, существует ли бакет, если нет - создаем его
	err = createBucketIfNotExists(client, config.Bucket)
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		Client:   client,
		Endpoint: config.Endpoint,
		Bucket:   config.Bucket,
	}, nil
}

func retryHeadBucket(client *s3.Client, bucket string, retries int, delay time.Duration) error {
	var err error
	for i := 0; i < retries; i++ {
		_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
			Bucket: &bucket,
		})
		if err == nil {
			return nil
		}

		// Логируем ошибку и ожидаем перед следующей попыткой
		log.Printf("Attempt %d: Error checking bucket: %v", i+1, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("failed to check bucket after %d attempts: %v", retries, err)
}

func createBucketIfNotExists(client *s3.Client, bucket string) error {
	err := retryHeadBucket(client, bucket, 5, 2*time.Second)
	if err != nil {
		log.Printf("Error checking bucket existence: %v", err)

		_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: &bucket,
		})
		if err != nil {
			return fmt.Errorf("unable to create bucket: %v", err)
		}
		log.Printf("Bucket %s created successfully.", bucket)

		err = applyBucketPolicy(client, bucket)
		if err != nil {
			return fmt.Errorf("failed to apply bucket policy: %v", err)
		}

		time.Sleep(5 * time.Second)

		err = uploadDefaultAvatar(client, bucket)
		if err != nil {
			return fmt.Errorf("failed to upload default avatar: %v", err)
		}
	}
	return nil
}

func applyBucketPolicy(client *s3.Client, bucket string) error {
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": ["s3:GetObject", "s3:PutObject"],
				"Resource": "arn:aws:s3:::` + bucket + `/*",
				"Principal": "*"
			}
		]
	}`

	_, err := client.PutBucketPolicy(context.TODO(), &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: &policy,
	})

	if err != nil {
		return fmt.Errorf("failed to apply bucket policy: %v", err)
	}

	log.Printf("Bucket policy applied to %s successfully.", bucket)
	return nil
}

func uploadDefaultAvatar(client *s3.Client, bucket string) error {
	defaultAvatar, err := os.ReadFile("default.webp")
	if err != nil {
		return fmt.Errorf("unable to read default avatar file: %v", err)
	}

	objectKey := "avatars/default.webp"
	uploadInput := &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &objectKey,
		Body:        bytes.NewReader(defaultAvatar),
		ContentType: aws.String("image/webp"),
	}

	_, err = client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		return fmt.Errorf("unable to upload default avatar to S3: %v", err)
	}

	log.Println("Default avatar uploaded successfully.")
	return nil
}
