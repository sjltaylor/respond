package respond

import (
	"fmt"
	"net/http"
)


////// Not Acceptable

type NotAcceptableError struct {
	Accept string
}

func NewNotAcceptableError(request *http.Request) *NotAcceptableError {
	return &NotAcceptableError{request.Header.Get(`Accept`)}
}

func (err *NotAcceptableError) Error() string {
	return fmt.Sprintf("unable to generate content satisfying Accept header: %s", err.Accept)
}

func (err *NotAcceptableError) HTTPStatusCode() int {
	return 406
}