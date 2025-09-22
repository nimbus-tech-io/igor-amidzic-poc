package s3

import (
	"fmt"
	"io"
	"time"

	"goapp/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	client *s3.S3
	bucket string
}

func NewS3Client() (*S3Client, error) {
	cfg := config.LoadConfig()
	
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.AWS.Region),
		Endpoint:         aws.String(cfg.AWS.Endpoint),
		S3ForcePathStyle: aws.Bool(cfg.S3.ForcePathStyle),
		Credentials:      credentials.NewStaticCredentials(cfg.AWS.Credentials.AccessKey, cfg.AWS.Credentials.SecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	client := s3.New(sess)

	return &S3Client{
		client: client,
		bucket: cfg.S3.BucketName,
	}, nil
}

func (s *S3Client) CreateBucketIfNotExists() error {
	_, err := s.client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err == nil {
		return nil 
	}

	_, err = s.client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	fmt.Printf("Created S3 bucket: %s\n", s.bucket)
	return nil
}

type UploadResult struct {
	Key      string `json:"key"`
	Location string `json:"location"`
	Size     int64  `json:"size"`
	Bucket   string `json:"bucket"`
}

func (s *S3Client) UploadFile(file io.Reader, filename string, contentType string, userID string) (*UploadResult, error) {
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("uploads/%s/%d_%s", userID, timestamp, filename)

	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        aws.ReadSeekCloser(file),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	headResult, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	var size int64
	if err == nil && headResult.ContentLength != nil {
		size = *headResult.ContentLength
	}

	return &UploadResult{
		Key:      key,
		Location: fmt.Sprintf("http://localhost:4566/%s/%s", s.bucket, key),
		Size:     size,
		Bucket:   s.bucket,
	}, nil
}

func (s *S3Client) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Client) GetFileURL(key string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return url, nil
}