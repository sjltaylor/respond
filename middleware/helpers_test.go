package middleware

import (
	"net/http"
)

type testEndpointWithMiddleware struct{}

func (ep *testEndpointWithMiddleware) Middlewares() []Middleware {

	var mw MiddlewareFunc = func(w http.ResponseWriter, r *http.Request, next NextFunc) error {
		if _, err := w.Write([]byte("HELLO ")); err != nil {
			panic(err)
		}
		return next(w)
	}

	return []Middleware{mw}
}

func (ep *testEndpointWithMiddleware) Process(w http.ResponseWriter, r *http.Request) error {
	if _, err := w.Write([]byte("WORLD")); err != nil {
		panic(err)
	}
	return nil
}
