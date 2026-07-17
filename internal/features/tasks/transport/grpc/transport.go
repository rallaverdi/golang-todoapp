package tasks_transport_grpc

import (
	"context"

	todov1 "github.com/rallaverdi/golang-todoapp/gen/go/todo/v1"
	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

type TasksService interface {
	GetTask(ctx context.Context, id int) (domain.Task, error)
}

type TasksGRPCHandler struct {
	todov1.UnimplementedTaskServiceServer

	tasksService TasksService
}

func NewTasksGRPCHandler(tasksService TasksService) *TasksGRPCHandler {
	return &TasksGRPCHandler{
		tasksService: tasksService,
	}
}

var _ todov1.TaskServiceServer = (*TasksGRPCHandler)(nil)
