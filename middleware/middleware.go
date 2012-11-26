package middleware

import "net/http"

type Middleware interface {
	Process(http.ResponseWriter, *http.Request, NextFunc) error
}

type Endpoint interface {
	Process(http.ResponseWriter, *http.Request) error
}

type NextFunc func(http.ResponseWriter) error
type MiddlewareFunc func(http.ResponseWriter, *http.Request, NextFunc) error
type EndpointFunc func(http.ResponseWriter, *http.Request) error

func (fn MiddlewareFunc) Process (response http.ResponseWriter, request *http.Request, next NextFunc) error {
	return fn(response, request, next)
}

func (fn EndpointFunc) Process (response http.ResponseWriter, request *http.Request) error {
	return fn(response, request)
}

type chain struct {
	middlewares []Middleware
}

func Chain(middlewares ...Middleware) *chain {
	return &chain{middlewares: middlewares}
}

func endpointAsMiddleware(ep Endpoint) Middleware {
	var mw MiddlewareFunc = func(response http.ResponseWriter, request *http.Request, next NextFunc) error {

		if err := ep.Process(response, request); err != nil {
			return err
		}

		return next(response)
	}

	return mw
}

func (c *chain) Chain(middlewares ...Middleware) *chain {

	newMiddlewares := append([]Middleware{}, c.middlewares...)
	newMiddlewares = append(newMiddlewares, middlewares...)
	return &chain{middlewares: newMiddlewares}
}

func (c *chain) EndpointFunc (fn EndpointFunc) http.HandlerFunc {
	var endpoint Endpoint
	endpoint = fn
	return c.Endpoint(endpoint)
}

func (c *chain) Endpoint(ep Endpoint) http.HandlerFunc {

	middlewares := append(c.middlewares, endpointAsMiddleware(ep))

	return func(response http.ResponseWriter, request *http.Request) {

		i := -1

		var current Middleware

		var next NextFunc

		next = func(w http.ResponseWriter) error {

			i++

			if i >= len(middlewares) {
				return nil
			}

			current = middlewares[i]

			return current.Process(w, request, next)
		}

		if err := next(response); err != nil {
			panic(err)
		}
	}
}
