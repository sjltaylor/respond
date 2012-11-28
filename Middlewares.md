# Respond: out of the box middleware

## Error Handling

This middleware is best placed as the first link in a middleware chain. If the request processing causes a panic or an error is returned from upstream server, the error handling middle delegates rendering to a user-defined function

Create an `ErrorHandlerMiddleware` and place it in a chain.

```
mw := middleware.Chain(respond.ErrorHandlerMiddleware(myErrorHandler))
```

`myErrorHandler` is of type `ServerErrorEndpointFunc`:

```
type ServerErrorEndpointFunc func(http.ResponseWriter, *http.Request, error) error
```

## Not Found middleware

## Logging middleware

