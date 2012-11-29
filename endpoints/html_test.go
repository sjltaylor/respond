package endpoints

import (
	"fmt"
	"testing"
	"net/http"
	"respond/testHelpers"
)

var failingHandler HTMLEndpointHandler = func (response http.ResponseWriter, request *http.Request) (interface{}, error) {
	return nil, fmt.Errorf("FAILED")
}

var successfulHandler HTMLEndpointHandler = func (response http.ResponseWriter, request *http.Request) (interface{}, error) {
	return "HELLO", nil
}

var exampleSuccessfulEndpoint *HTMLEndpoint = NewHTMLEndpoint("test-layout", "one", "nested/two").Handler(successfulHandler)
var exampleFailingEndpoint    *HTMLEndpoint = NewHTMLEndpoint("test-layout", "one", "nested/two").Handler(failingHandler)

var exampleFailingRenderEndpoint *HTMLEndpoint = NewHTMLEndpoint("DOES NOT EXIST", "nor do i", "nested/or me!").Handler(successfulHandler)

func init () {
	TemplatesDirectory = "./test_templates"
}

func TestHTMLEndpointRendersTemplateWithTheDataReturnedByTheHandler (t *testing.T) {

	response := testHelpers.NewFakeResponseWriter()
	
	if err := exampleSuccessfulEndpoint.Process(response, nil); err != nil {
		t.Fatalf("processing failed: %s", err)
	}

	expectedBodyContent := `Hello. This is one. And this is two.` 
	
	if response.Body != expectedBodyContent {
		t.Fatalf("body content not rendered correctly. got '%s' instead of '%s'.", response.Body, expectedBodyContent)
	}

	if response.Header().Get(`Content-Type`) != `text/html` {
		t.Fatalf("response content type should be text/html")
	}
}

func TestHTMLEndpointReturnsAnErrorIfRenderingFails (t *testing.T) {

	if err := exampleFailingRenderEndpoint.Process(nil, nil); err == nil {
		t.Fatal("should return an error when rendering fails")
	}
}

func TestHTMLEndpointReturnsAnErrorIfTheHandlerReturnsAnError (t *testing.T) {

	if err := exampleFailingEndpoint.Process(nil, nil); err == nil {
		t.Fatalf("should return an error if the handler fails")
	}
}