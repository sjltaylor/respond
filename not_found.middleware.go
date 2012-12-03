package respond

import (
	"net/http"
	"respond/middleware"
)

type NotFoundErrorEndpointFunc func(http.ResponseWriter, *http.Request, *NotFoundError) error

type NotFoundMiddleware struct {
	NotFoundEndpoint NotFoundErrorEndpointFunc
}

func NewNotFoundMiddleware(errorPageEndpoint NotFoundErrorEndpointFunc) *NotFoundMiddleware {
	return &NotFoundMiddleware{
		NotFoundEndpoint: errorPageEndpoint,
	}
}

func (notFound *NotFoundMiddleware) Process(response http.ResponseWriter, request *http.Request,
	next middleware.NextFunc) (returnError error) {

	if err := next(response); err != nil {

		if notFoundError, ok := err.(*NotFoundError); ok {
			
			response.WriteHeader(404)

			return notFound.NotFoundEndpoint(response, request, notFoundError)
		}

		return err
	}

	return nil
}
