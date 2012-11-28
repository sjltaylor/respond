package endpoint

import (
	"net/http"
	"html/template"
)

var ReloadTemplates bool
var TemplatesDirectory string = "./webapp/html"

type HTMLEndpoint struct {
	Layout   string // **/<Layout>.layout.tmpl, from the root of TemplatesDirectory
	Partials []string // **/<?>.tmpl from the root of the TemplatesDirectory
	Handler func (http.ResponseWriter, *http.Request) interface{}, error
	cachedTemplate *template.Template
}

func NewHTMLEndpoint (layout string, partials ...string) *HTMLEndpoint {
	
	return &HTMLEndpoint{
		Layout: layout,
		Partials: partials,
	}
}

func (endpoint *HTMLEndpoint) Handler (fn func (http.ResponseWriter, *http.Request) error) *HTMLEndpoint {
	endpoint.Handler = fn
	return endpoint
}

func (endpoint *HTMLEndpoint) template () *template.Template {
	
	var t *template.Template

	if ReloadTemplates {

	}

}

func (endpoint *HTMLEndpoint) Process (response http.ResponseWriter, request *http.Request) error {
	// if the Handler returns an error this should return that error
	// if the request does not accept text/html -> 406

}