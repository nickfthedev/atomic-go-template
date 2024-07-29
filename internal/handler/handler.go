package handler

import (
	"my-go-template/internal/database"
)

type Handler struct {
	db database.Service
}

func NewHandler(db database.Service) *Handler {
	return &Handler{db: db}
}
