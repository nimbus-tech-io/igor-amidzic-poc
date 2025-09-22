package handlers

import (
	"encoding/json"
	"goapp/internal/models"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.RegisterResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}
	
	if req.Email == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.RegisterResponse{
			Success: false,
			Message: "Email and password are required",
		})
		return
	}
	
	user, err := h.repo.CreateUser(r.Context(), req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.RegisterResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.RegisterResponse{
		Success: true,
		Message: "User created successfully",
		User:    user.ToResponse(),
	})
}
	