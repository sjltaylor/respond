package respond

import (
	"fmt"
	"net/http"
)

////// Not Found

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

func (err *NotFoundError) HTTPStatusCode() int {
	return 404
}


////// Not Acceptable

type NotAcceptableError struct {
	Accept string
}

func NewNotAcceptableError(request *http.Request) *NotAcceptableError {
	return &NotAcceptableError{request.Header.Get(`Accept`)}
}

func (err NotAcceptableError) Error() string {
	return fmt.Sprintf("unable to generate content satisfying Accept header: %s", err.Accept)
}

func (err *NotAcceptableError) HTTPStatusCode() int {
	return 406
}


////// Bad Request

type BadRequestError struct {
	Detail string
}

func NewBadRequestError(format string, data ...interface{}) *BadRequestError {
	return &BadRequestError{Detail: fmt.Sprintf(format, data)}
}

func (err BadRequestError) Error() string {
	return err.Detail
}

func (err *BadRequestError) HTTPStatusCode() int {
	return 400
}
