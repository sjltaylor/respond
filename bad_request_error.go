package respond

import "fmt"


////// Bad Request

type BadRequestError struct {
	Detail string
}

func NewBadRequestError(format string, data ...interface{}) *BadRequestError {
	return &BadRequestError{Detail: fmt.Sprintf(format, data)}
}

func (err *BadRequestError) Error() string {
	return err.Detail
}

func (err *BadRequestError) HTTPStatusCode() int {
	return 400
}
