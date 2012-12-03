package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONEndpoint struct {
	handler Handler
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
	endpoint.handler = fn
	return endpoint
}

func (endpoint *JSONEndpoint) Process(response http.ResponseWriter, request *http.Request) (returnError error) {

	defer func() {

		if err := recover(); err != nil {
			returnError = fmt.Errorf("json endpoint: render failed: %s", err)
		}
	}()

	data, err := endpoint.handler(response, request)

	if err != nil {
		panic(err)
	}

	var payload []byte

	if payload, err = json.Marshal(data); err != nil {
		panic(err)
	}

	if _, err = response.Write(payload); err != nil {
		panic(err)
	}

	response.Header().Add(`Content-Type`, `application/json`)

	return nil
}
