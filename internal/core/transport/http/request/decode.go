package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return fmt.Errorf("decode json: %v:%w", err, core_errors.ErrInvalidArgument)
	}

	if err := requestValidator.Struct(dto); err != nil {
		return fmt.Errorf("request validation %v: %w", err, core_errors.ErrInvalidArgument)
	}
	return nil
}
