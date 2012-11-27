package respond

import (
	"fmt"
	"net/http"
	"testing"
	"respond/middleware"
)

func testErrorPageEndpoint(called *bool) ServerErrorEndpointFunc {

	return func(response http.ResponseWriter, request *http.Request, err error) error {
		*called = true
		return nil
	}
}

func TestErrorHandlerMiddlewareHandlesReturnedErrors(t *testing.T) {

	var called bool
	ehm := NewErrorHandlerMiddleware(testErrorPageEndpoint(&called))

	var nextFunc middleware.NextFunc = func(response http.ResponseWriter) error {
		return fmt.Errorf("WALLOP")
	}

	rw := newFakeResponseWriter()

	ehm.Process(rw, nil, nextFunc)

	if rw.status != 500 {
		t.Fatalf("response status should be 500 but was %d", rw.status)
	}

	if !called {
		t.Fatal("error page should have been called")
	}
}

func TestErrorHandlerMiddlewareReturnsAnErrorIfTheErrorPageEndpointReturnsAnError(t *testing.T) {

	ehm := NewErrorHandlerMiddleware(func(response http.ResponseWriter, request *http.Request, err error) error {
		return fmt.Errorf("ERROR RENDERING ERROR PAGE (oh NO!)")
	})

	var nextFunc middleware.NextFunc = func(response http.ResponseWriter) error {
		panic("ERROR!")
	}

	rw := newFakeResponseWriter()

	e := ehm.Process(rw, nil, nextFunc)

	if rw.status != 500 {
		t.Fatalf("response status should be 500 but was %d", rw.status)
	}

	if e == nil {
		t.Fatal("should return any error returned by the error page handler")
	}
}

func TestErrorHandlerMiddlewareHandlesPanics(t *testing.T) {

	var called bool
	ehm := NewErrorHandlerMiddleware(testErrorPageEndpoint(&called))

	var nextFunc middleware.NextFunc = func(response http.ResponseWriter) error {
		panic("BANG")
	}

	rw := newFakeResponseWriter()

	ehm.Process(rw, nil, nextFunc)

	if rw.status != 500 {
		t.Fatalf("response status should be 500 but was %d", rw.status)
	}

	if !called {
		t.Fatal("error page should have been called")
	}
}
