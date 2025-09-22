package queue

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	client   *sqs.SQS
	queueURL string
}

type FileStats struct {
	UserID       string `json:"user_id"`
	UserEmail    string `json:"user_email"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	S3Key        string `json:"s3_key"`
	S3Bucket     string `json:"s3_bucket"`
	ContentType  string `json:"content_type"`
	UploadedAt   string `json:"uploaded_at"`
}

func NewSQSClient() (*SQSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Endpoint:    aws.String("http://localhost:4566"), 
		Credentials: credentials.NewStaticCredentials("test", "test", ""), 
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	client := sqs.New(sess)
	queueName := "nimbus-queue"

	queueURL, err := getOrCreateQueue(client, queueName)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue URL: %v", err)
	}

	return &SQSClient{
		client:   client,
		queueURL: queueURL,
	}, nil
}

func getOrCreateQueue(client *sqs.SQS, queueName string) (string, error) {
	result, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err == nil {
		return *result.QueueUrl, nil
	}

	createResult, err := client.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}

	fmt.Printf("Created SQS queue: %s\n", queueName)
	return *createResult.QueueUrl, nil
}

func (s *SQSClient) SendFileStats(stats FileStats) error {
	messageBody, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal file stats: %v", err)
	}

	_, err = s.client.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String("file_upload"),
			},
			"user_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(stats.UserID),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send message to queue: %v", err)
	}

	fmt.Printf("Sent file stats to queue: %s\n", stats.FileName)
	return nil
}