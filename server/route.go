package server

import (
	"net/http"

	"github.com/pkg/errors"
)

// Route is a Handler matched by a regex
// it routes to handlers with different methods
// its ErrorHandler will be called if no handler matches with the requested method
type Route struct {
	regex        string
	methods      map[string]http.Handler
	errorHandler ErrorHandler
}

// NewRoute creates a new Router with the given regex and ErrorHandler, and returns its pointer
// If nil ErrorHandler is provided it will fall back to the BasicErrorHandler
func NewRoute(pathRegex string, errHandler ErrorHandler) *Route {
	var handler ErrorHandler = &BasicErrorHandler{}
	if errHandler != nil {
		handler = errHandler
	}
	return &Route{
		regex:        pathRegex,
		methods:      make(map[string]http.Handler),
		errorHandler: handler,
	}
}

// WithMethod adds the given Handler with the given method to the Route and returns it
func (route *Route) WithMethod(method string, handler http.Handler) *Route {
	route.methods[method] = handler
	return route
}

// AddMethod adds the given Handler with the given method to the Route
func (route *Route) AddMethod(method string, handler http.Handler) {
	route.methods[method] = handler
}

// ServeHTTP
// The Route will try to match the request's method with its know methods and
// pass the request accordingly, if no matching method is find it will call
// the ErrorHandler with HTTP 405 MethodNotAllowed
func (route *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for method, handler := range route.methods {
		if req.Method == method {
			handler.ServeHTTP(res, req)
			return
		}
	}
	route.errorHandler.HandleError(res, req, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
}
