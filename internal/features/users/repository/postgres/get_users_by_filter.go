package users_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsersByFilter(ctx context.Context, filter domain.UsersFilter) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT id, version, full_name, phone_number FROM todoapp.users`)
	args := []any{}
	conditions := []string{}

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf(`id="$%d"`, len(args)+1))
		args = append(args, *filter.UserID)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	args = append(args, filter.PageSize)
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)))

	args = append(args, filter.PageNumber)
	queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)))

	var userModels []UserModel

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf(`select users: %w`, err)
	}
	defer rows.Close()
	for rows.Next() {
		var userModel UserModel
		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users error: %w", err)
		}
		userModels = append(userModels, userModel)
	}
	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil

}
