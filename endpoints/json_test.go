package endpoints

import (
	"fmt"
	"net/http"
	"respond/test_helpers"
	"testing"
)

type JSONEndpointTest struct {
	successful     *JSONEndpoint
	failingHandler *JSONEndpoint
	failingRender  *JSONEndpoint
}

var jsonEndpointTest *JSONEndpointTest

func init() {

	jsonEndpointTest = &JSONEndpointTest{}

	var failingHandler Handler = func(response http.ResponseWriter, request *http.Request) (interface{}, error) {

		return nil, fmt.Errorf("FAILED")
	}

	var successfulHandler Handler = func(response http.ResponseWriter, request *http.Request) (interface{}, error) {

		data := make(map[string]interface{})

		data["field1"] = 123
		data["field2"] = []string{"one", "two", "three"}

		return data, nil
	}

	var handlerReturningUnserializableJSON Handler = func(response http.ResponseWriter, request *http.Request) (interface{}, error) {
		return make(chan string), nil
	}

	jsonEndpointTest.successful = NewJSONEndpoint().Handler(successfulHandler)
	jsonEndpointTest.failingHandler = NewJSONEndpoint().Handler(failingHandler)
	jsonEndpointTest.failingRender = NewJSONEndpoint().Handler(handlerReturningUnserializableJSON)
}

func TestJSONEndpointRendersTemplateWithTheDataReturnedByTheHandler(t *testing.T) {

	response := testHelpers.NewFakeResponseWriter()

	if err := jsonEndpointTest.successful.Process(response, nil); err != nil {
		t.Fatalf("processing failed: %s", err)
	}

	expectedBodyContent := `{"field1":123,"field2":["one","two","three"]}`

	if response.Body != expectedBodyContent {
		t.Fatalf("data was not correctly encoded as JSON. got '%s' instead of '%s'.", response.Body, expectedBodyContent)
	}

	if response.Header().Get(`Content-Type`) != `application/json` {
		t.Fatalf("response content type should be application/json")
	}
}

func TestJSONEndpointReturnsAnErrorIfRenderingFails(t *testing.T) {

	if err := jsonEndpointTest.failingRender.Process(nil, nil); err == nil {
		t.Fatal("should return an error when rendering fails")
	}
}

func TestJSONEndpointReturnsAnErrorIfTheHandlerReturnsAnError(t *testing.T) {

	if err := jsonEndpointTest.failingHandler.Process(nil, nil); err == nil {
		t.Fatalf("should return an error if the handler fails")
	}
}
