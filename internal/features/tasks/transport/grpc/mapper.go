package tasks_transport_grpc

import (
	todov1 "github.com/rallaverdi/golang-todoapp/gen/go/todo/v1"
	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func taskProtoFromDomain(task domain.Task) *todov1.Task {
	response := &todov1.Task{
		Id:           int64(task.ID),
		Version:      int64(task.Version),
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    timestamppb.New(task.CreatedAt),
		AuthorUserId: int64(task.AuthorUserID),
	}

	if task.CompletedAt != nil {
		response.CompletedAt = timestamppb.New(*task.CompletedAt)
	}

	return response
}
