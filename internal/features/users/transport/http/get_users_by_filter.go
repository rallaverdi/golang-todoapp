package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) GetUsersByFilter(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	limit, offset, userID, err := getLimitOffsetUserIDQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit' or 'offset' query param",
		)
		return
	}

	if limit != nil && offset != nil {
		filters := domain.NewUsersFilter(*limit, *offset, userID)
		userDomains, err := h.usersService.GetUsersByFilter(ctx, filters)
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

}

func getLimitOffsetUserIDQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
		userIDQueryParamKey = "user_id"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset query param: %w", err)
	}

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return offset, limit, nil, nil
	}
	return limit, offset, userID, nil
}
