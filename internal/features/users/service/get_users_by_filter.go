package users_service

import (
	"context"
	"fmt"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

func (s *UsersService) GetUsersByFilter(ctx context.Context, filter domain.UsersFilter) ([]domain.User, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf(`invalid filter args:%w`, core_errors.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsersByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get users by filter: %w", err)
	}

	return users, nil
}
