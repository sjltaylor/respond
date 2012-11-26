package respond

import (
	"fmt"
	"net/http"
)

type NotFoundError struct {
	ResourceDescription string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", err.ResourceDescription)
}

func NewNotFoundError(format string, data ...interface{}) *NotFoundError {
	return &NotFoundError{
		ResourceDescription: fmt.Sprintf(format, data...),
	}
}

type NotAcceptableError struct {
	Accept string
}

func NewNotAcceptableError(request *http.Request) *NotAcceptableError {
	return &NotAcceptableError{request.Header.Get(`Accept`)}
}
