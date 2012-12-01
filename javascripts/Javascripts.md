# Javascripts

* uses requirejs
* assumes we are using http://www.gorillatoolkit.org/pkg/mux
* renders the javascript or returns an error incl. a respond.NotFoundError
* the routh *must* match the path of the javascript to a patter called "javascript". e.g.: router.HandleFunc("/js/{javascript:.*}", ...