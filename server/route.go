package server

import (
	"net/http"

	"github.com/pkg/errors"
)

type Route struct {
	regex        string
	methods      map[string]http.Handler
	errorHandler ErrorHandler
}

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

func (route *Route) WithMethod(method string, handler http.Handler) *Route {
	route.methods[method] = handler
	return route
}

func (route *Route) AddMethod(method string, handler http.Handler) {
	route.methods[method] = handler
}

func (route *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for method, handler := range route.methods {
		if req.Method == method {
			handler.ServeHTTP(res, req)
			return
		}
	}
	route.errorHandler.HandleError(res, req, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
}
