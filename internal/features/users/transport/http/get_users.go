package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers			godoc
// @Summary 		Get user list
// @Description		Get user list with optional pagination
// @Tags 			users
// @Produce 		json
// @Param 			limit query int false "Page size of users"
// @Param 			offset query int false "Page number of users"
// @Success			200 {object} GetUsersResponse "Get users successfully"
// @Failure			400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure			500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router			/users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit' or 'offset' query param",
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'users' data",
		)
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)

	if err != nil {
		return nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get offset query param: %w", err)
	}
	return limit, offset, nil
}
