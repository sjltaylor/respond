package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"respond/testHelpers"
)

func TestMiddlewareErrorDefault(t *testing.T) {

	/*

		a http handler that wraps middleware/endpoints which return an error should panic

	*/

	middlewares := Middlewares()

	httpHandler := middlewares.EndpointFunc(func(w http.ResponseWriter, r *http.Request) error {

		return fmt.Errorf("ERROR!")
	})

	defer func() {

		e, ok := recover().(error)

		if !ok || e.Error() != "ERROR!" {
			t.Fatalf("the error from the middlewares was passed to panic: %s", e)
		}
	}()

	// this should cause a panic
	httpHandler(nil, &http.Request{})
}

func TestMiddlewareComposition(t *testing.T) {

	var mw1 MiddlewareFunc

	mw1 = func(w http.ResponseWriter, r *http.Request, next NextFunc) error { return nil }

	middlewares := Middlewares(mw1)

	func () {
		expectedCount := 1
		actualCount := len(middlewares)
		
		if actualCount != expectedCount {
			t.Fatalf("wrong number of middlewares: %d instead of %d", actualCount, expectedCount)
		}
	}()

	extendedMiddlewares := middlewares.And(mw1)

	func () {
		expectedCount := 2
		actualCount := len(extendedMiddlewares)
		
		if actualCount != expectedCount {
			t.Fatalf("wrong number of middlewares in extendedMiddlewares: %d instead of %d", actualCount, expectedCount)
		}
	}()
}

func TestMiddlewareChainExecution(t *testing.T) {

	response := testHelpers.NewFakeResponseWriter()
	request := &http.Request{}

	var lastMiddlewareCalled byte

	// middleware 1 
	var mw1Response http.ResponseWriter
	var mw1Request *http.Request
	var mw1ResponseWrapper http.ResponseWriter = testHelpers.NewFakeResponseWriter()
	var mw1StashedError error

	var mw1 MiddlewareFunc = func(w http.ResponseWriter, r *http.Request, next NextFunc) error {
		mw1Response = w
		mw1Request = r
		lastMiddlewareCalled = 1
		mw1StashedError = next(mw1ResponseWrapper)
		return nil
	}

	var mw2Response http.ResponseWriter
	var mw2Request *http.Request

	var mw2 MiddlewareFunc = func(w http.ResponseWriter, r *http.Request, next NextFunc) error {

		mw2Response = w
		mw2Request = r
		lastMiddlewareCalled = 2

		return next(w)
	}

	var endpointCalled bool
	endpointError := fmt.Errorf("OH NO! something went wrong in the endpoint function")

	httpHandler := Middlewares(mw1, mw2).EndpointFunc(func(w http.ResponseWriter, r *http.Request) error {
		endpointCalled = true
		return endpointError
	})

	// call the http handler
	httpHandler(response, request)

	if !endpointCalled {
		t.Fatal(`the endpoint was not called`)
	}

	if lastMiddlewareCalled != 2 {
		t.Fatal(`middlewares must be called in the order of arguments to Chain(...)`)
	}

	if mw1Response != response {
		t.Fatal(`the middleware function was not passed the response`)
	}

	if mw2Response != mw1ResponseWrapper {
		t.Fatal(`the 'next NextFunc' does not propogate the response passed to it by the middleware function`)
	}

	if mw2Request != request && mw1Request != request {
		t.Fatal(`the request is not propogated through the middle ware`)
	}

	if mw1StashedError != endpointError {
		t.Fatal(`errors are not returned through the middleware`)
	}
}
