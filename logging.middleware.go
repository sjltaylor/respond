package respond

import(
	"log"
	"net/http"
	"respond/middleware"
)

type LoggingMiddleware struct {}

func (loggingMiddleware *LoggingMiddleware) Process(response http.ResponseWriter, request *http.Request, next middleware.NextFunc) error {
	
	log.Printf("respond: %s request to: %s", request.Method, request.URL.Path)
	
	err := next(response)

	if err != nil {
		log.Printf("respond: error: %s", err)
	}

	return err
}