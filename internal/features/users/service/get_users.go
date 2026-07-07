package users_service

import (
	"context"
	"fmt"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

func (s *UsersService) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(`limit must be a non-negative integer:%w`, core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(`offset must be a non-negative integer:%w`, core_errors.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting users list: %w", err)
	}
	return users, nil

}
