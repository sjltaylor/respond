package respond

import (
	"fmt"
	"respond/middleware"
	"testing"
	"net/http"
)

func TestNotFoundErrorMiddleware (t *testing.T) {

	called := false

	var notFoundEndpoint NotFoundErrorEndpointFunc = func (response http.ResponseWriter, request *http.Request, err *NotFoundError) error {
		called = true
		return nil
	}

	notFoundMiddleware := NewNotFoundMiddleware(notFoundEndpoint)

	var next middleware.NextFunc = func (response http.ResponseWriter) error {
		return NewNotFoundError("something important")
	}

	e := notFoundMiddleware.Process(nil, nil, next)

	if !called {
		t.Fatal(`NotFoundErrorEndpoint should have been called`)
	}

	if e != nil {
		t.Fatal(`a NotFoundError should have not have been returned`)
	}

	called = false

	next = func (response http.ResponseWriter) error {
		return fmt.Errorf("some other error")
	}

	e = notFoundMiddleware.Process(nil, nil, next)

	if called {
		t.Fatal("should not call NotFoundErrorEndpoint when the error is not a NotFoundError")
	}

	if e == nil {
		t.Fatal("error should have been return for the next middle up to handle")
	}
}

