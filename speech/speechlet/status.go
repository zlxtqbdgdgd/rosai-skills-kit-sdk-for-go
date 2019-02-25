package speechlet

import "errors"

type Reason string

const (

	// The user explicitly exited the skill
	USER_INITIATED Reason = "USER_INITIATED"
	// An error occurred and the session had to be ended
	ERROR Reason = "ERROR"
	// The rosai skill has not received a valid response within the maximum allowed
	// number of re-prompts.
	EXCEEDED_MAX_REPROMPTS Reason = "EXCEEDED_MAX_REPROMPTS"
)

func NewGoodStatus() *Status {
	return &Status{}
}

func NewInternalErrStatus(detail string) *Status {
	return &Status{
		Code:         ApiInternal,
		ErrorDetails: detail,
	}
}

func NewMismatchStatus(detail string) *Status {
	return &Status{
		Code:         ApiServiceMismatched,
		ErrorDetails: detail,
	}
}

type Status struct {
	Code         ApiStatusCode `json:"code"`
	ErrorType    string        `json:"errorType,omitempty"`
	ErrorDetails string        `json:"errorDetails,omitempty"`
}

type ApiStatusCode int

var (
	ErrServiceMismatched = errors.New("service_mismatched")
	ErrServiceInternal   = errors.New("service_internal_error")
)

const (
	ApiSuccess  ApiStatusCode = 0
	ApiNoResult ApiStatusCode = 1
	// 客户端错误状态码
	ApiBadRequest   ApiStatusCode = 400
	ApiUnauthorized ApiStatusCode = 401
	//ApiNotFound     ApiStatusCode = 404
	// 服务端错误状态码
	ApiInternal     ApiStatusCode = 500
	ApiNotSupported ApiStatusCode = 501
	// 第三方服务错误状态码
	ApiServiceUnavailable   ApiStatusCode = 601
	ApiServiceUnknownFormat ApiStatusCode = 602
	ApiServiceMismatched    ApiStatusCode = 603
)
