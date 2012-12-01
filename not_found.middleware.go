package respond

import (
	"respond/middleware"
	"net/http"
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

func (notFound *NotFoundMiddleware) Process (response http.ResponseWriter, request *http.Request,
	next middleware.NextFunc) (returnError error) {

	if err := next(response); err != nil {
		if notFoundError, ok := err.(*NotFoundError); ok {
			return notFound.NotFoundEndpoint(response, request, notFoundError)
		}	

		return err
	}

	return nil
} 