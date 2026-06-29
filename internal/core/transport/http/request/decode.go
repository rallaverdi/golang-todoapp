package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}

	if err := requestValidator.Struct(dto); err != nil {
		return fmt.Errorf("error validating request body: %w", err)
	}
	return nil
}
