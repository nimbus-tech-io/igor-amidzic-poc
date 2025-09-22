package handlers

import (
	"goapp/internal/queue"
	"goapp/internal/repo"
	"goapp/internal/s3"
)

type Handler struct {
	repo      repo.Repository
	s3Client  *s3.S3Client
	sqsClient *queue.SQSClient
}

func NewHandler(repository repo.Repository, s3Client *s3.S3Client, sqsClient *queue.SQSClient) *Handler {
	return &Handler{
		repo:      repository,
		s3Client:  s3Client,
		sqsClient: sqsClient,
	}
}