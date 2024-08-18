package services

import (
	"context"
	"vk_tarantool_project/internal/config"
	"vk_tarantool_project/internal/domain"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, username string) (*domain.UserInfo, error)
}

type DataRepository interface {
	WriteData(ctx context.Context, data domain.Data) error
	// ReadData(ctx context.Context, keys domain.DataKeys) (domain.Data, error)
}

type Service struct {
	urep UserRepository
	drep DataRepository
	conf *config.Config
}

// New create new service to work with UserRepository and DataRepository
func New(conf *config.Config, urep UserRepository, drep DataRepository) *Service {
	return &Service{
		urep: urep,
		drep: drep,
		conf: conf,
	}
}

func (s *Service) WriteData(ctx context.Context, data domain.Data) error {

	if len(data.Data) == 0 {
		return nil
	}

	return s.drep.WriteData(ctx, data)
}

func (s *Service) ReadData(ctx context.Context, keys domain.DataKeys) (domain.Data, error) {
	//TODO implement me
	panic("implement me")
}
