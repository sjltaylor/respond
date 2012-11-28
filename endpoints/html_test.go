package endpoints

import (
	"fmt"
	"testing"
	"net/http"
)

failingHandler := func (response http.ResponseWriter, request *http.Request) interface{}, error {
	return nil, fmt.Errorf("FAILED")
}

successfulHandler := func (response http.ResponseWriter, request *http.Request) interface{}, error {
	return "HELLO", nil
}


exampleSuccessfulEndpoint := NewHTMLEndpoint("test-layout", "one", "two", "three").Handler(successfulHandler)
exampleFailingEndpoint    := NewHTMLEndpoint("test-layout", "one", "two", "three").Handler(failingHandler)

func TestHTMLEndpoint (t *testing.T) {

	

	//ep := 
	
	// should render 406 if the request does not accept html
	// should render the template with the data returned by the handler
	// should return an error and not render anything if there is an error
}