package services

import (
	"context"
	"vk_tarantool_project/internal/config"
	"vk_tarantool_project/internal/domain"
	"vk_tarantool_project/internal/pkg/jwt"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, username string) (*domain.UserInfo, error)
}

// TODO: Сделать указатели
type DataRepository interface {
	WriteData(ctx context.Context, data domain.Data) error
	ReadData(ctx context.Context, keys domain.DataKeys) (domain.Data, error)
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

func (s *Service) Login(ctx context.Context, user domain.UserInfo) (string, error) {

	userDB, err := s.urep.GetUserByName(ctx, user.Username)
	if err != nil {
		return "", err
	}

	if userDB.Password != user.Password {
		return "", domain.ErrPasswordMismatch
	}

	token, err := jwt.NewToken(user, s.conf.Secret, s.conf.TokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) WriteData(ctx context.Context, data domain.Data) error {
	if len(data.Data) == 0 {
		return nil
	}
	return s.drep.WriteData(ctx, data)
}

func (s *Service) ReadData(ctx context.Context, keys domain.DataKeys) (domain.Data, error) {
	data, err := s.drep.ReadData(ctx, keys)
	return data, err
}
