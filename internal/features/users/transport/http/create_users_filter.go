package users_transport_http

import (
	"net/http"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
)

type CreateUsersFilterRequest struct {
	PageSize   int  `json:"page_size" validate:"required"`
	PageNumber int  `json:"page_number" validate:"required"`
	UserID     *int `json:"user_id" validate:"omitempty"`
}

type CreateUsersFilterResponse struct {
	FilterID string `json:"filter_id"`
}

func (h *UsersHTTPHandler) CreateUserFilter(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	var request CreateUsersFilterRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "invalid users filter")
		return
	}

	filters := domain.NewUsersFilter(request.PageSize, request.PageNumber, request.UserID)

	filterID, err := h.usersService.CreateUsersFilter(ctx, filters)
	if err != nil {
		responseHandler.ErrorResponse(err, "create users filter")
		return
	}

	response := CreateUsersFilterResponse{FilterID: filterID}
	responseHandler.JSONResponse(response, http.StatusCreated)

}
