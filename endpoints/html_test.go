package endpoints

import (
	"fmt"
	"net/http"
	"respond/test_helpers"
	"testing"
	"respond"
	"respond/middleware"
)

type HTMLEndpointTest struct {
	success         *HTMLEndpoint
	failingHandling *HTMLEndpoint
	failingRender   *HTMLEndpoint
}

var htmlEndpointTest *HTMLEndpointTest

func init() {
	TemplatesDirectory = "./test_templates"
	htmlEndpointTest = buildHtmlEndpointTest()
}

func buildHtmlEndpointTest() *HTMLEndpointTest {

	htmlEndpointTest := &HTMLEndpointTest{}

	var failingHandler Handler = func(response http.ResponseWriter, request *http.Request) (interface{}, error) {
		return nil, fmt.Errorf("FAILED")
	}

	var successfulHandler Handler = func(response http.ResponseWriter, request *http.Request) (interface{}, error) {
		return "HELLO", nil
	}

	htmlEndpointTest.success = NewHTMLEndpoint("test-layout", "one", "nested/two").Handler(successfulHandler)
	htmlEndpointTest.failingHandling = NewHTMLEndpoint("test-layout", "one", "nested/two").Handler(failingHandler)

	htmlEndpointTest.failingRender = NewHTMLEndpoint("DOES NOT EXIST", "nor do i", "nested/or me!").Handler(successfulHandler)

	return htmlEndpointTest
}

func TestHTMLEndpointRendersTemplateWithTheDataReturnedByTheHandler(t *testing.T) {

	response := testHelpers.NewFakeResponseWriter()

	if err := htmlEndpointTest.success.Process(response, nil); err != nil {
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

func TestHTMLEndpointMiddlewareInterceptsUnacceptableRequests(t *testing.T) {

	request, _ := http.NewRequest("GET", "http://localhost", nil)
	request.Header.Add(`Accept`, `image/jpeg`)

	var err error

	var errorIntercept middleware.MiddlewareFunc = func (w http.ResponseWriter, r *http.Request, next middleware.NextFunc) error {
		err = next(w)
		return nil
	}

	httpHandler := middleware.Middlewares(errorIntercept).Endpoint(htmlEndpointTest.success)

	httpHandler.ServeHTTP(nil, request)
	
	if _, ok := err.(*respond.NotAcceptableError); !ok {
		t.Fatalf("a NotAcceptableError was not return. got: %s", err)
	}
}

func TestHTMLEndpointReturnsAnErrorIfRenderingFails(t *testing.T) {

	if err := htmlEndpointTest.failingRender.Process(nil, nil); err == nil {
		t.Fatal("should return an error when rendering fails")
	}
}

func TestHTMLEndpointReturnsAnErrorIfTheHandlerReturnsAnError(t *testing.T) {

	if err := htmlEndpointTest.failingHandling.Process(nil, nil); err == nil {
		t.Fatalf("should return an error if the handler fails")
	}
}
