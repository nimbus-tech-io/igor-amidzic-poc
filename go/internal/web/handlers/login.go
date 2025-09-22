package handlers

import (
	"encoding/json"
	"goapp/internal/auth"
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
    
    // Generate JWT token
    jwtToken, err := auth.GenerateJWT(user.ID, user.Email)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(models.LoginResponse{
            Success: false,
            Message: "Failed to generate token",
        })
        return
    }
    
    // Set JWT cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "jwt_token",
        Value:    jwtToken,
        Path:     "/",
        MaxAge:   86400, // 24 hours
        HttpOnly: true,
        Secure:   false, // Set to true in production with HTTPS
        SameSite: http.SameSiteLaxMode,
    })
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(models.LoginResponse{
        Success: true,
        Message: "Login successful",
        User: user.ToResponse(),
    })
}