package handlers

import (
	"encoding/json"
	"goapp/internal/auth"
	"goapp/internal/queue"
	"log"
	"net/http"
	"time"
)

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to parse form",
		})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No file provided",
		})
		return
	}
	defer file.Close()

	result, err := h.s3Client.UploadFile(file, header.Filename, header.Header.Get("Content-Type"), user.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to upload file to S3",
			"error":   err.Error(),
		})
		return
	}

	fileStats := queue.FileStats{
		UserID:      user.UserID,
		UserEmail:   user.Email,
		FileName:    header.Filename,
		FileSize:    result.Size,
		S3Key:       result.Key,
		S3Bucket:    result.Bucket,
		ContentType: header.Header.Get("Content-Type"),
		UploadedAt:  time.Now().Format(time.RFC3339),
	}

	if err := h.sqsClient.SendFileStats(fileStats); err != nil {
		log.Printf("Failed to send file stats to queue: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"message":      "File uploaded successfully to S3",
		"key":          result.Key,
		"location":     result.Location,
		"bucket":       result.Bucket,
		"originalName": header.Filename,
		"size":         header.Size,
		"uploadedBy":   user.Email,
	})
}