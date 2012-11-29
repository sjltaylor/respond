package endpoints

import (
	"net/http"
	"path"
	"html/template"
	"fmt"
	"bytes"
)

var ReloadTemplates bool
var TemplatesDirectory string = "./webapp/html"

type HTMLEndpointHandler func (http.ResponseWriter, *http.Request) (interface{}, error)

type HTMLEndpoint struct {
	Layout   string // **/<Layout>.layout.tmpl, from the root of TemplatesDirectory
	Partials []string // **/<?>.tmpl from the root of the TemplatesDirectory
	handler HTMLEndpointHandler
	cachedTemplate *template.Template
}

func NewHTMLEndpoint (layout string, partials ...string) *HTMLEndpoint {
	
	return &HTMLEndpoint{
		Layout: layout,
		Partials: partials,
	}
}

func (endpoint *HTMLEndpoint) Handler (fn HTMLEndpointHandler) *HTMLEndpoint {
	endpoint.handler = fn
	return endpoint
}

func (endpoint *HTMLEndpoint) loadTemplates () (t *template.Template, err error) {

	// stdlib template libraries require the main template to be named after the first file in its set
	
	layoutFilename := endpoint.Layout + `.layout.tmpl`
	
	t = template.New(layoutFilename)

	filepaths := []string{path.Join(TemplatesDirectory, layoutFilename)}

	for _, partial := range endpoint.Partials {
		filepaths = append(filepaths, path.Join(TemplatesDirectory, partial + `.tmpl`))
	}

	t, err = t.ParseFiles(filepaths...)

	if err != nil {
		return nil, fmt.Errorf("parse failed: %+v, %s", endpoint, err)
	}

	return
} 

func (endpoint *HTMLEndpoint) template () (*template.Template, error) {
	
	if ReloadTemplates {
		return endpoint.loadTemplates()
	}

	if endpoint.cachedTemplate == nil {
		var err error
		
		if endpoint.cachedTemplate, err = endpoint.loadTemplates(); err != nil {
			return nil, err
		}
	}

	return endpoint.cachedTemplate, nil
}

func (endpoint *HTMLEndpoint) Process (response http.ResponseWriter, request *http.Request) (returnError error) {
	
	defer func() {

		if err := recover(); err != nil {
			returnError = fmt.Errorf("html endpoint: render failed: %s", err)
		}
	}()

	response.Header().Add(`Content-Type`, `text/html`)

	t, err := endpoint.template()

	if err != nil {
		panic(err)
	}

	data, err := endpoint.handler(response, request)

	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBufferString("")

	if err = t.Execute(buffer, data); err != nil {
		panic(err)
	}

	if _, err = response.Write([]byte(buffer.String())); err != nil {
		panic(err)
	}

	return nil
}


