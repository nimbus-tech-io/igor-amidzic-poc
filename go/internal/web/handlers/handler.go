package handlers

import (
	"goapp/internal/repo"
)

type Handler struct {
	repo repo.Repository
}

func NewHandler(repository repo.Repository) *Handler {
	return &Handler{repo: repository}
}