package handlers

import (
	"encoding/json"
	"goapp/internal/models"
	"net/http"
)



func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.LoginResponse{
            Success: false,
            Message: "Invalid JSON format",
        })
        return
    }
    
    if req.Email == "" || req.Password == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.LoginResponse{
            Success: false,
            Message: "Email and password are required",
        })
        return
    }
    
    user, err := h.repo.VerifyPassword(r.Context(), req.Email, req.Password)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(models.LoginResponse{
            Success: false,
            Message: "Invalid email or password",
        })
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(models.LoginResponse{
        Success: true,
        Message: "Login successful",
        User: user.ToResponse(),
    })
}