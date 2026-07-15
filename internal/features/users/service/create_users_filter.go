package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

func (s *UsersService) CreateUsersFilter(ctx context.Context, filter domain.UsersFilter) (string, error) {

	if err := filter.Validate(); err != nil {
		return "", fmt.Errorf("validate filter: %w", err)
	}

	filterID, err := s.usersCache.FindFilterID(ctx, filter)
	if err == nil {
		return filterID, nil
	}

	if !errors.Is(err, core_errors.ErrNotFound) {
		return "", fmt.Errorf("find cached filter error: %w", err)
	}

	users, err := s.usersRepository.GetUsersByFilter(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("get users by filter: %w", err)
	}

	filterID = uuid.NewString()
	if err := s.usersCache.SaveResult(ctx, filter, filterID, users); err != nil {
		return "", fmt.Errorf("save filter to cache result error: %w", err)
	}

	return filterID, nil
}
