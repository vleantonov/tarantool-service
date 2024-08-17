package services

import (
	"context"
	"vk_tarantool_project/internal/config"
	"vk_tarantool_project/internal/domain"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, username string) (*domain.UserInfo, error)
}

type Service struct {
	urep UserRepository
	conf *config.Config
}

// New create new service to work with UserRepository and DataRepository
func New(conf *config.Config, urep UserRepository) *Service {
	return &Service{
		urep: urep,
		conf: conf,
	}
}
