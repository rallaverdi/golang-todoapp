package tasks_transport_grpc

import (
	"context"

	todov1 "github.com/rallaverdi/golang-todoapp/gen/go/todo/v1"
	core_grpc_transport "github.com/rallaverdi/golang-todoapp/internal/core/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *TasksGRPCHandler) GetTask(
	ctx context.Context,
	request *todov1.GetTaskRequest,
) (*todov1.Task, error) {
	if request == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	id := request.GetId()
	if id <= 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"id must be greater than zero",
		)
	}

	taskID := int(id)
	if int64(taskID) != id {
		return nil, status.Error(
			codes.InvalidArgument,
			"id is outside the supported range",
		)
	}

	task, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		return nil, core_grpc_transport.ToStatusError(err)
	}

	return taskProtoFromDomain(task), nil
}
