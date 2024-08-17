package handlers

import (
	"context"
	"vk_tarantool_project/internal/domain"
)

type Service interface {
	Login(ctx context.Context, info domain.UserInfo) (string, error)
}

type Handler struct {
	service Service
}

// New Create a new http app handler with using service
func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}