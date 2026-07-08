package tasks_service

import (
	"context"
	"fmt"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

func (s *TasksService) GetTasks(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(`offset must be a non-negative integer:%w`, core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(`offset must be a non-negative integer:%w`, core_errors.ErrInvalidArgument)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting tasks list from repository: %w", err)
	}
	return tasks, nil
}
