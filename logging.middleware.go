	package respond

import (
	"log"
	"net/http"
	"respond/middleware"
)

type LoggingMiddleware struct{}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (loggingMiddleware *LoggingMiddleware) Process(response http.ResponseWriter, request *http.Request, next middleware.NextFunc) error {

	log.Printf("respond: %s %s", request.Method, request.URL.Path)

	err := next(response)

	if err != nil {
		log.Printf("respond: error from upstream request handlers: %s", err)
	}

	return err
}
