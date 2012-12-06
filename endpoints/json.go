package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"respond"
	"respond/middleware"
)

type JSONHandler func(http.ResponseWriter, *http.Request) ([]byte, error)

type JSONEndpoint struct {
	handler JSONHandler
}

func emptyJsonHandler(respond http.ResponseWriter, request *http.Request) (interface{}, error) {
	
	return (map[string]interface{}{}), nil
}

func NewJSONEndpoint() *JSONEndpoint {

	ep := &JSONEndpoint{}
	ep.Handler(emptyJsonHandler)

	return ep
}

func (endpoint *JSONEndpoint) Handler(fn Handler) *JSONEndpoint {
	
	return endpoint.JSONHandler(func (response http.ResponseWriter, request *http.Request) ([]byte, error) {
		
		data, err := fn(response, request)

		if err != nil {
			return nil, err
		}

		var payload []byte

		if payload, err = json.Marshal(data); err != nil {
			return nil, err
		}
		
		return payload, err		
	})
}

func (endpoint *JSONEndpoint) JSONHandler(fn JSONHandler) *JSONEndpoint {
	endpoint.handler = fn
	return endpoint
}

func (endpoint *JSONEndpoint) Middlewares() []middleware.Middleware {
	return []middleware.Middleware{respond.NewAcceptFilterMiddleware(`application/json`)}
}

func (endpoint *JSONEndpoint) Process(response http.ResponseWriter, request *http.Request) (returnError error) {

	defer func() {

		if err := recover(); err != nil {
			returnError = fmt.Errorf("json endpoint: render failed: %s", err)
		}
	}()

	payload, err := endpoint.handler(response, request)

	if err != nil {
		panic(err)
	}

	if _, err = response.Write(payload); err != nil {
		panic(err)
	}

	response.Header().Add(`Content-Type`, `application/json`)

	return nil
}
