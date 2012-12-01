package endpoints

import "net/http"

type Handler func (http.ResponseWriter, *http.Request) (interface{}, error)
