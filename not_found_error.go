package respond

import "fmt"

////// Not Found

type NotFoundError struct {
	ResourceDescription string
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", err.ResourceDescription)
}

func NewNotFoundError(format string, data ...interface{}) *NotFoundError {
	return &NotFoundError{
		ResourceDescription: fmt.Sprintf(format, data...),
	}
}

func (err *NotFoundError) HTTPStatusCode() int {
	return 404
}
