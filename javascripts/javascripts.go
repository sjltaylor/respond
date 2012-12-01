package javascripts

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"respond"
)

type JavascriptEndpoint struct {
	ScriptPath string
}

func (ep *JavascriptEndpoint) Process(response http.ResponseWriter, req *http.Request) (err error) {

	script := path.Join(ep.ScriptPath, mux.Vars(req)["javascript"])

	var file *os.File

	if file, err = os.Open(script); err != nil {
		return respond.NewNotFoundError("cannot open javascript file: %s, %s", script, err)
	}
	defer file.Close()

	var fileInfo os.FileInfo

	if fileInfo, err = file.Stat(); err != nil {
		return fmt.Errorf("can't stat javascript file: %s, %s", script, err)
	} else if fileInfo.IsDir() {
		return respond.NewNotFoundError("cannot open javascript: %s is a directory", script)
	}

	var bytes []byte

	if bytes, err = ioutil.ReadAll(file); err != nil {
		return fmt.Errorf("can't read javascript file: %s, %s", script, err)
	}

	response.Header().Add("Content-Type", "application/javascript")
	response.Write(bytes)

	return nil
}
