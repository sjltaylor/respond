package respond

import (
	"fmt"
	"log"
	"net/http"
	"respond/middleware"
)

type ServerErrorEndpointFunc func(http.ResponseWriter, *http.Request, error) error

type StatusCodedError interface {
	HTTPStatusCode() int
}

type ErrorHandlerMiddleware struct {
	ServerErrorEndpoint ServerErrorEndpointFunc
}

func NewErrorHandlerMiddleware(errorPageEndpoint ServerErrorEndpointFunc) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		ServerErrorEndpoint: errorPageEndpoint,
	}
}

func (errorHandler *ErrorHandlerMiddleware) Process(response http.ResponseWriter, request *http.Request,
	next middleware.NextFunc) (returnError error) {

	defer func() {

		if recovered := recover(); recovered != nil {

			var err error

			if e, ok := recovered.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%+v", recovered)
			}

			log.Println(fmt.Sprintf("respond: unhandled error: %s", err))

			if statusCodedError, ok := err.(StatusCodedError); ok {
				response.WriteHeader(statusCodedError.HTTPStatusCode())
			} else {
				response.WriteHeader(500)
			}

			err = errorHandler.ServerErrorEndpoint(response, request, err)

			if err != nil {
				returnError = fmt.Errorf("%s", err)
			}
		}
	}()

	if err := next(response); err != nil {
		panic(err)
	}

	return nil
}
