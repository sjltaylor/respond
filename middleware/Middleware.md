# Respond-Middleware

Middleware system for web daemons in Go. Depends only on `net/http`.

## Why

The Go standard library facilitates http request handling with http.HandleFunc

```
func (http.ResponseWriter, *http.Request)
```

Web daemons must deal with:
 
* Error handling
* Authentication
* Logging
* Caching
* etc...

These concerns pervade the logic required to handle most if not all of the http endpoints for a program.

Middleware is a way to encapsulate these concerns, greatly simplifying the logic at each endpoint.


## What

Middleware is modelled as a sequence of 0..n functions followed by an endpoint function. 

This sequence of function is enclosed in a `http.HandlerFunc` which can be passed to `http.HandleFunc(…)`.

Note: other libraries such as [Gorilla web tooklit](http://www.gorillatoolkit.org/pkg/mux) also use `http.HandlerFunc`.


## How


### Define a middleware chain

```
import "respond/middleware"
chain := middleware.Chain(mw1, mw2, …)
```
`mw1` and `mw2` implement `middleware.Middleware`

```
type Middleware interface {
	Process(http.ResponseWriter, *http.Request, NextFunc) error
}
```
For convenience middleware also defines...
```
type MiddlewareFunc func(http.ResponseWriter, *http.Request, NextFunc) error
```
which implements the `Middleware` interface. Thus, middleware can be define just as a function:

```
var mw middleware.Middleware

mw = func (response http.ResponseWriter, request *http.Request, next NextFunc) error {
	// middleware implementation goes here
	return nil
}

middleware.Chain(mw, …)
```

### Example middleware function

```
func NoOpMiddleware (
	response http.ResponseWriter, 
	request *http.Request, 
	next NextFunc) error {
	
	/*
		parameter 'next' is used to call the next function in the 
		middleware chain.
		
		NextFunc takes a response writer as a parameter which 
		allows middlewares to wrap the response writer before passing it on.
	*/
	
	err := next(response)
	
	/*
		err is the error returned from middlewares/endpoints further 
		up the chain.
	*/
	
	return err
	
}
```

### Create an endpoint

An endpoint has the signature: 

```
type EndpointFunc func(http.ResponseWriter, *http.Request) error
```

For example:

```
func renderSomething (response http.ResponseWriter, request *http.Request) error {
	
	json, err := getSomeJsonFromSomewhere()
	
	// assume that there is some error handling middleware...
	if err != nil {
		return err
	}

	response.Write([]byte(json))
	return nil
}

```
The error is returned and handled by error handling middleware further down the chain.

To create a `http.HandlerFunc` from this endpoint using middleware we call `Endpoint` or `EndpointFunc` on the chain:

```
http.HandleFunc("/something", chain.EndpointFunc(endpointFunc))
```
or if we have something which implements `middleware.Endpoint`:

```
http.HandleFunc("/something", chain.Endpoint(implementorOfEndpoint))
```

### Chains from chains

```
longerChain := chain.Chain(mw3, mw4, …)
```

`longerChain` contains all of the middleware in `chain` plus `mw3` and `mw4`. It is equivalent to `Chain(mw1, mw2, mw3, mw4)`


## Limitations

* Only works with `net/http`


