package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Armen Lipitskiy"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79032234356"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'FullName' can not be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("'FullName' length must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("'PhoneNumber' length must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("'PhoneNumber' value must start with '+'")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser		godoc
// @Summary 		Update user data
// @Description		Update data of existing user
// @Description		### Tree State Logic:
// @Description		1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description		2. **Явно передано значение**: `"phone_number"`: `"+79992223322"`- устанавливает новый номер телефона в БД
// @Description		3. **Передан null**: `"phone_number"`: `"null"` - очищает поле в БД (set to NULL)
// @Description		4. **Ограничения**: `full_name` не может быть выставлен как null
// @Tags 			users
// @Accept			json
// @Produce 		json
// @Param 			id path int true "ID of user to update"
// @Param 			request body PatchUserRequest true "Patch user request body"
// @Success			200 {object} PatchUserResponse "User has been updated successfully"
// @Failure			400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure			404 {object} core_http_response.ErrorResponse "User not found"
// @Failure			409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure			500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router			/users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	userID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)
	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fields:\nFullName: %v\nPhoneNumber: %v",
			request.FullName,
			request.PhoneNumber,
		),
	)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
