package server

import (
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

type Router struct {
	routes       []*Route
	errorHandler ErrorHandler
}

func NewRouter(errHandler ErrorHandler) *Router {
	var handler ErrorHandler = &BasicErrorHandler{}
	if errHandler != nil {
		handler = errHandler
	}
	return &Router{
		errorHandler: handler,
	}
}

func (router *Router) WithRoute(route *Route) *Router {
	router.routes = append(router.routes, route)
	return router
}

func (router *Router) WithRoutes(routes ...*Route) *Router {
	router.routes = append(router.routes, routes...)
	return router
}

func (router *Router) AddRoute(route *Route) {
	router.routes = append(router.routes, route)
}

func (router *Router) AddRoutes(routes ...*Route) {
	router.routes = append(router.routes, routes...)
}

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

func urlMatch(url, regex string) bool {
	matched, err := regexp.MatchString(regex, url)
	if err != nil {
		return false
	}
	return matched
}
