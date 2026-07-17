package core_grpc_transport

import (
	"context"
	"errors"

	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatusError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, core_errors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, core_errors.ErrConflict):
		return status.Error(codes.Aborted, err.Error())

	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")

	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
