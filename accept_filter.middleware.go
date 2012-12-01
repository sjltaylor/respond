package respond

import (
	"net/http"
	"respond/accept"
	"respond/middleware"
)

type AcceptFilterMiddleware struct {
	ContentType string
}

func NewAcceptFilterMiddleware(contentType string) *AcceptFilterMiddleware {
	return &AcceptFilterMiddleware{ContentType: contentType}
}

func (acceptFilter *AcceptFilterMiddleware) Process(response http.ResponseWriter, request *http.Request,
	next middleware.NextFunc) error {

	acceptHeader, err := accept.ParseAcceptHeader(request)

	if err != nil {
		return NewBadRequestError(err.Error())
	}

	if acceptHeader.AcceptsMediaType(acceptFilter.ContentType) {

		return next(response)
	}

	return NewNotAcceptableError(request)
}
