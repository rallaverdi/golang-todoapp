package users_transport_http

import (
	"net/http"

	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) GetCachedUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	filterID := getFilterIDQueryParam(r)
	if filterID == "" {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "Invalid filter ID")
		return
	}

	users, err := h.usersService.GetCachedUsers(ctx, filterID)
	if err != nil {
		responseHandler.ErrorResponse(err, "get users by filter id")
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(users))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func getFilterIDQueryParam(r *http.Request) string {
	const (
		filterIDQueryParamKey = "limit"
	)
	return r.PathValue(filterIDQueryParamKey)
}
