package users_service

import (
	"context"
	"fmt"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

func (s *UsersService) GetCachedUsers(ctx context.Context, filterID string) ([]domain.User, error) {

	if filterID == "" {
		return []domain.User{}, fmt.Errorf("filterID is empty but required: %w", core_errors.ErrInvalidArgument)
	}

	users, err := s.usersCache.GetUsers(ctx, filterID)
	if err != nil {
		return []domain.User{}, fmt.Errorf("get users from cache error: %w", err)
	}
	return users, nil
}
