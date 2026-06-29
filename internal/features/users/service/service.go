package users_service

import (
	"context"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{}
}
