package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		errorText := fmt.Errorf("param='%s' by key='%s' is not a valid integer: %v:%w", param, key, err, core_errors.ErrInvalidArgument)
		return nil, errorText
	}
	return &val, nil
}
