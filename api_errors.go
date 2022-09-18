package exchange

import (
	"fmt"
)

// NewApiError creates a new ApiError
func NewApiError(msg string) *ApiError {
	return &ApiError{msg: msg}
}

// ApiError type for ExchangeRate-API response
type ApiError struct {
	msg string
}

// ApiError message
func (e *ApiError) Error() string {
	return e.msg
}

// ErrType error type in ExchangeRate-API request
type ErrType string

const (
	UnsupportedCode     ErrType = "unsupported-code"
	MalformedRequest    ErrType = "malformed-request"
	InvalidKey          ErrType = "invalid-key"
	InactiveAccount     ErrType = "inactive-account"
	QuotaReached        ErrType = "quota-reached"
	PlanUpgradeRequired ErrType = "plan-upgrade-required"
)

func responseApiError(r Response) error {
	if r.Result == SUCCESS {
		return nil
	}

	switch r.ErrType {
	case UnsupportedCode:
		return NewApiError("Unsupported code")
	case MalformedRequest:
		return NewApiError("Malformed request")
	case InvalidKey:
		return NewApiError("Invalid key")
	case InactiveAccount:
		return NewApiError("Inactive account")
	case QuotaReached:
		return NewApiError("Quota reached")
	case PlanUpgradeRequired:
		return NewApiError("Plan upgrade required")
	default:
		return NewApiError(fmt.Sprintf("Unexpected error type: %s", r.ErrType))
	}
}
