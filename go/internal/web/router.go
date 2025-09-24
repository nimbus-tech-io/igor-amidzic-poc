package web

import (
	"goapp/internal/auth"
	"goapp/internal/queue"
	"goapp/internal/repo"
	"goapp/internal/s3"
	"goapp/internal/web/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		
		next.ServeHTTP(w, r)
	})
}

func NewRouter(repository repo.Repository, s3Client *s3.S3Client, sqsClient *queue.SQSClient) *mux.Router {
	r := mux.NewRouter()
	
	r.Use(corsMiddleware)
	
	handler := handlers.NewHandler(repository, s3Client, sqsClient)
	authMiddleware := auth.NewMiddleware()

	r.HandleFunc("/auth/register", handler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth/logout", handler.Logout).Methods("POST", "OPTIONS")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	protected.HandleFunc("/upload", handler.UploadFile).Methods("POST")
	
	return r
}