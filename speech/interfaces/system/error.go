package system

type ErrorType string

// Represents the types of errors permissible
// INVALID_RESPONSE
// The request was invalid. E.g. malformed, unauthorized, forbidden, not found, etc
// DEVICE_COMMUNICATION_ERROR
// A problem occurred communicating with the device
// INTERNAL_SERVICE_ERROR
// The server was unable to process the request as expected
const (
	INVALID_RESPONSE           ErrorType = "INVALID_RESPONSE"
	DEVICE_COMMUNICATION_ERROR ErrorType = "DEVICE_COMMUNICATION_ERROR"
	INTERNAL_SERVICE_ERROR     ErrorType = "INTERNAL_SERVICE_ERROR"
)

func NewError(typ ErrorType, message string) Error {
	return Error{typ, message}
}

type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return "error type: " + string(e.Type) + ", message: " + e.Message
}
