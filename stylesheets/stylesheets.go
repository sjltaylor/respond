package stylesheets

import (
	"fmt"
	"os/exec"
	"path/filepath"
	// "io/ioutil"
	//  "path"
	//  "bytes"
	//  "html/template"
	"code.google.com/p/gorilla/mux"
	"net/http"
)

type StylesheetEndpoint struct {
	StylesheetPath string
}

func (ep *StylesheetEndpoint) Process(response http.ResponseWriter, request *http.Request) error {

	stylesheet := filepath.Join(ep.StylesheetPath, mux.Vars(request)["stylesheet"])

	cmd := exec.Command("lessc", stylesheet)

	output, err := cmd.CombinedOutput()

	// TODO: check for the existence of the file and return respond.NotFoundError
	if err != nil {
		return fmt.Errorf("rendering stylesheet failed: %s, %s\n%s", stylesheet, err, string(output))
	}

	response.Header().Add(`Content-Type`, `text/css`)
	_, err = response.Write(output)

	return err
}
