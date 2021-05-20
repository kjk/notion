package notion

import (
	"encoding/json"
	"errors"
	"fmt"
)

// See: https://developers.notion.com/reference/errors.
var (
	ErrInvalidJSON        = errors.New("notion: request body could not be decoded as JSON")
	ErrInvalidRequestURL  = errors.New("notion: request URL is not valid")
	ErrInvalidRequest     = errors.New("notion: request is not supported")
	ErrValidation         = errors.New("notion: request body does not match the schema for the expected parameters")
	ErrUnauthorized       = errors.New("notion: bearer token is not valid")
	ErrRestrictedResource = errors.New("notion: client doesn't have permission to perform this operation")
	ErrObjectNotFound     = errors.New("notion: the resource does not exist")
	ErrConflict           = errors.New("notion: the transaction could not be completed, potentially due to a data collision")
	ErrRateLimited        = errors.New("notion: this request exceeds the number of requests allowed")
	ErrInternalServer     = errors.New("notion: an unexpected error occurred")
	ErrServiceUnavailable = errors.New("notion: service is unavailable")
)

var errMap = map[string]error{
	"invalid_json":          ErrInvalidJSON,
	"invalid_request_url":   ErrInvalidRequestURL,
	"invalid_request":       ErrInvalidRequest,
	"validation_error":      ErrValidation,
	"unauthorized":          ErrUnauthorized,
	"restricted_resource":   ErrRestrictedResource,
	"object_not_found":      ErrObjectNotFound,
	"conflict_error":        ErrConflict,
	"rate_limited":          ErrRateLimited,
	"internal_server_error": ErrInternalServer,
	"service_unavailable":   ErrServiceUnavailable,
}

type APIError struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements `error`.
func (err *APIError) Error() string {
	return fmt.Sprintf("%v (code: %v, status: %v)", err.Message, err.Code, err.Status)
}

func (err *APIError) Unwrap() error {
	mapped, ok := errMap[err.Code]
	if !ok {
		return fmt.Errorf("notion: %v", err.Error())
	}

	return mapped
}

func parseErrorResponseJSON(d []byte) error {
	var apiErr APIError

	err := json.Unmarshal(d, &apiErr)
	if err != nil {
		return fmt.Errorf("failed to parse error from HTTP response: %w", err)
	}

	return &apiErr
}
