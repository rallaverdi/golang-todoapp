package tasks_service

import (
	"context"
	"fmt"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository error: %w", err)
	}
	return task, nil
}
