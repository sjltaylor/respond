# Stylesheets

* uses lessc
* assumes we are using http://www.gorillatoolkit.org/pkg/mux
* renders the stylesheet or returns an error (TODO: incl. a respond.NotFoundError)
* the routh *must* match the path of the stylesheet to a pattern called "stylesheet". e.g.: router.HandleFunc("/css/{stylesheet:.*}", ...