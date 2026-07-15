package users_service

import (
	"context"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
	usersCache      UsersRedisCache
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error)
	GetUsersByFilter(ctx context.Context, filters domain.UsersFilter) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error)
}

type UsersRedisCache interface {
	FindFilterID(ctx context.Context, filter domain.UsersFilter) (string, error)
	SaveResult(
		ctx context.Context,
		filter domain.UsersFilter,
		filterID string,
		users []domain.User,
	) error
	GetUsers(ctx context.Context, filterID string) ([]domain.User, error)
}

func NewUsersService(usersRepository UsersRepository, usersCache UsersRedisCache) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		usersCache:      usersCache,
	}
}
