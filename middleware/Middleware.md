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

HTTP requests are received by the first middleware which decided whether to forward the request 'upstream' to the next middleware, ultimately reach an endpoint. 

Endpoints and middleware can return an error to be handled by downstream middlewares.

## How


### Define a middleware sequence

```
import "respond/middleware"
middlewares := middleware.Middlewares(mw1, mw2, …)
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

middleware.Middlewares(mw, …)
```

### Example middleware function

```
func NoOpMiddleware (
	response http.ResponseWriter, 
	request *http.Request, 
	next NextFunc) error {
	
	/*
		parameter 'next' is used to indicate that the request should be passed onto the next middleware.
		
		NextFunc takes a response writer as a parameter which 
		allows middlewares to wrap the response writer before passing it on.
	*/
	
	err := next(response)
	
	/*
		err is the error returned from middlewares/endpoints further 
		upstream.
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
http.HandleFunc("/something", middlewares.EndpointFunc(endpointFunc))
```
or if we have something which implements `middleware.Endpoint`:

```
http.HandleFunc("/something", middlewares.Endpoint(implementorOfEndpoint))
```

### Middlewares And...

```
extendedMiddleware := middlewares.And(mw3, mw4, …)
```

`extendedMiddleware` contains all of `middlewares` and `mw3` and `mw4`. It is equivalent to `Middlewares(mw1, mw2, mw3, mw4)`

## Endpoints With Middleware

when passing an endpoint to `Endpoint()`, if the endpoint implements `EndpointWithMiddlewares` it's middlewares are appended to the middlewares:

```
type EndpointWithMiddleware interface {
	Middlewares() []Middleware
}
```

For example, in respond/endpoints, the `HTMLEndpoint` defines middlewares containing an accept filter middleware that filters out requests that do not accept media type text/html and returns a NotAcceptableError backdown the middlewares

## Limitations

* Only works with `net/http`


