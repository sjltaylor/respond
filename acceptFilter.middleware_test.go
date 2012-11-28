package respond

import (
	"fmt"
	"testing"
	"respond/middleware"
	"net/http"
	"reflect"
)

func TestAcceptFilterMiddlewareFiltersRequestsThatCannotBeFullfilled (t *testing.T) {

	afmw := NewAcceptFilterMiddleware(`application/json`)

	request, _ := http.NewRequest(`GET`, `http://somedomain`, nil)
	request.Header.Add(`Accept`, `text/html`)

	var called bool

	var next middleware.NextFunc = func(_ http.ResponseWriter) error {
		called = true
		return nil
	}

	e := afmw.Process(nil, request, next)

	if called {
		t.Fatal(`the request should have been filtered`)
	}

	if _, ok := e.(*NotAcceptableError); !ok {
		t.Fatalf("should return a NotAcceptable error, but returned: %s: %s", reflect.TypeOf(e), e)
	}
}

func TestAcceptFilterMiddlewareDoesNotFilterRequestsThatCanBeFullfilled (t *testing.T) {

	afmw := NewAcceptFilterMiddleware(`application/json`)

	request, _ := http.NewRequest(`GET`, `http://somedomain`, nil)
	request.Header.Add(`Accept`, `application/json`)

	var called bool

	var upstreamError error = fmt.Errorf("UPSTREAM ERROR")

	var next middleware.NextFunc = func(_ http.ResponseWriter) error {
		called = true
		return upstreamError
	}

	e := afmw.Process(nil, request, next)

	if !called {
		t.Fatal(`the request should not have been filtered`)
	}

	if e != upstreamError {
		t.Fatal(`should return any upstream error`)
	}
}

func TestAcceptFilterMiddlewareReturnsABadRequestErrorIfTheAcceptHeaderIsMalformed (t *testing.T) {

	afmw := NewAcceptFilterMiddleware(`application/json`)

	request, _ := http.NewRequest(`GET`, `http://somedomain`, nil)
	request.Header.Add(`Accept`, `not a syntactically valid media range`)

	var called bool

	var next middleware.NextFunc = func(_ http.ResponseWriter) error {
		called = true
		return nil
	}

	e := afmw.Process(nil, request, next)

	if called {
		t.Fatalf("the request should have been filtered. Filter Accept")
	}

	if _, ok := e.(*BadRequestError); !ok {
		t.Fatalf("should return a BadRequestError, but returned: %s: %s", reflect.TypeOf(e), e)
	}
}

// test it returns a not acceptable error
// test if is a noop if the Accept head can be fullfilled
// test it returns a 400 if the Accept header can't be parsed