package server

import (
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

// Router is a Handler with a list of Routes and an ErrorHandler
// It tries to match the request's URL with its routes using the route's regex
// If a match found the request will be passed to the route
// If no match found the ErrorHandler will be called with HTTP 404 Not Found
//
// As the Router checks its routes in order it is best to declare them:
// from more specific ones to less specific ones
type Router struct {
	routes       []*Route
	errorHandler ErrorHandler
}

// NewRouter creates a new Router with the given ErrorHandler and return its pointer
func NewRouter(errHandler ErrorHandler) *Router {
	var handler ErrorHandler = &BasicErrorHandler{}
	if errHandler != nil {
		handler = errHandler
	}
	return &Router{
		errorHandler: handler,
	}
}

// WithRoute adds a Route to the Router and returns it
func (router *Router) WithRoute(route *Route) *Router {
	router.routes = append(router.routes, route)
	return router
}

// WithRoutes adds all Routes to the Router and returns it
func (router *Router) WithRoutes(routes ...*Route) *Router {
	router.routes = append(router.routes, routes...)
	return router
}

// AddRoute adds a Route to the Router
func (router *Router) AddRoute(route *Route) {
	router.routes = append(router.routes, route)
}

// AddRoutes adds all the Routes to the Router
func (router *Router) AddRoutes(routes ...*Route) {
	router.routes = append(router.routes, routes...)
}

// ServeHTTP
// The Router will try to match the request's URL with its Route's regex
// If a match it will pass the request to the matching Route
// If no match found the Router's ErrorHandler will be called with HTTP 404 Not Found
func (router *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range router.routes {
		if urlMatch(req.URL.String(), route.regex) {
			route.ServeHTTP(res, req)
			return
		}
	}
	router.errorHandler.HandleError(
		res,
		req,
		http.StatusNotFound,
		errors.Errorf("Not found: %s", req.URL.String()))
}

// urlMatch is a helper function that matches the given string against the given regex
// It will return true if it matches, and false if the check fails or it does not matches
// TODO: the fail case could be an internal server error
func urlMatch(url, regex string) bool {
	matched, err := regexp.MatchString(regex, url)
	if err != nil {
		return false
	}
	return matched
}
