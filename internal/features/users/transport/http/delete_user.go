package users_transport_http

import (
	"net/http"

	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/request"

	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
)

// DeleteUser 		godoc
// @Summary 		Delete user
// @Description		Delete user from system by ID
// @Tags 			users
// @Param			id path int true "ID of user to delete"
// @Success			204 "Successfully deleted"
// @Failure			400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure			404 {object} core_http_response.ErrorResponse "User not found"
// @Failure			500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router			/users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed getting user id",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()

}
