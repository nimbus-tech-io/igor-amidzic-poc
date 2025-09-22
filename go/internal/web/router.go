package web

import (
	"goapp/internal/auth"
	"goapp/internal/queue"
	"goapp/internal/repo"
	"goapp/internal/s3"
	"goapp/internal/web/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(repository repo.Repository, s3Client *s3.S3Client, sqsClient *queue.SQSClient) *mux.Router {
	r := mux.NewRouter()
	
	handler := handlers.NewHandler(repository, s3Client, sqsClient)
	authMiddleware := auth.NewMiddleware()
	
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("POST")
	
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	protected.HandleFunc("/upload", handler.UploadFile).Methods("POST")
	
	return r
}