package main

import (
	"context"
	"goapp/internal/provider"
	"goapp/internal/web"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	
	appProvider := provider.NewAppProvider()
	repo := appProvider.Repository()
	s3Client := appProvider.S3Client()
	sqsClient := appProvider.SQSClient()
	
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable not set")
	}
	
	if err := repo.Connect(ctx, dsn); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer repo.Close()
	
	r := web.NewRouter(repo, s3Client, sqsClient)
	
	log.Println("Server starting on :8080")
	log.Println("LocalStack S3 endpoint: http://localhost:4566")
	log.Fatal(http.ListenAndServe(":8080", r))
}