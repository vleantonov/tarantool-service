package services

import (
	"context"
	"vk_tarantool_project/internal/domain"
	"vk_tarantool_project/internal/pkg/jwt"
)

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
