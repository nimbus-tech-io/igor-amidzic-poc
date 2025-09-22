package provider

import (
	"goapp/internal/queue"
	"goapp/internal/repo"
	"goapp/internal/s3"
	"log"
)

type AppProvider struct {
    repository repo.Repository
    s3Client   *s3.S3Client
    sqsClient  *queue.SQSClient
}

func NewAppProvider() *AppProvider {
    return &AppProvider{}
}

func (p *AppProvider) Repository() repo.Repository {
    if p.repository == nil {
        p.repository = repo.NewNeonRepo()
    }
    return p.repository
}

func (p *AppProvider) S3Client() *s3.S3Client {
    if p.s3Client == nil {
        var err error
        p.s3Client, err = s3.NewS3Client()
        if err != nil {
            log.Fatal("Failed to initialize S3 client:", err)
        }
        
        // Create bucket if it doesn't exist
        if err := p.s3Client.CreateBucketIfNotExists(); err != nil {
            log.Fatal("Failed to create S3 bucket:", err)
        }
    }
    return p.s3Client
}

func (p *AppProvider) SQSClient() *queue.SQSClient {
    if p.sqsClient == nil {
        var err error
        p.sqsClient, err = queue.NewSQSClient()
        if err != nil {
            log.Fatal("Failed to initialize SQS client:", err)
        }
    }
    return p.sqsClient
}
