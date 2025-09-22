package web

import (
	"goapp/internal/repo"
	"goapp/internal/web/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(repository repo.Repository) *mux.Router {
	r := mux.NewRouter()
	
	handler := handlers.NewHandler(repository)
	
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	
	return r
}