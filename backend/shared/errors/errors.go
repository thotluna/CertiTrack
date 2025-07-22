package errors

import "fmt"

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("status: %d, message: %s", e.Status, e.Message)
}

func NewErrorResponse(status int, message string) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: message,
	}
}

var (
	ErrInternalServer = NewErrorResponse(500, "Internal server error")
	ErrNotFound       = NewErrorResponse(404, "Resource not found")
	ErrBadRequest     = NewErrorResponse(400, "Bad request")
)
