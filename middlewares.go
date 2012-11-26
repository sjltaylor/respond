package respond

import (
	"fmt"
	"net/http"
	"respond/middleware"
)

type ErrorPageEndpoint func(response http.ResponseWriter, request *http.Request, err error) error

type ErrorHandlerMiddleware struct {
	ErrorPage ErrorPageEndpoint
}

func NewErrorHandlerMiddleware(errorPageEndpoint ErrorPageEndpoint) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		ErrorPage: errorPageEndpoint,
	}
}

func (errorHandler *ErrorHandlerMiddleware) Process (response http.ResponseWriter, request *http.Request,
	next middleware.NextFunc) (returnError error) {

	defer func() {

		if recovered := recover(); recovered != nil {

			response.WriteHeader(500)

			var err error
			
			if e, ok := recovered.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%+v", e)
			}			

			if errorHandler.ErrorPage == nil {
				_, err = response.Write([]byte(`Server Error`))
			} else {
				err = errorHandler.ErrorPage(response, request, err)
			}

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
