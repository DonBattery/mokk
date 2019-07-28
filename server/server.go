package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Server is a wrapper for httptest.Server
type Server struct {
	*httptest.Server
}

// NewServer creates a new HTTP TestServer with the supplyed handler
// It is immediately initialized and can be reached at Server.URL
func NewServer(router http.Handler) *Server {
	return &Server{
		Server: httptest.NewServer(router),
	}
}

// TestServer is a wrapper for Server created in a testing context
type TestServer struct {
	*Server
	router *Router
	test   *testing.T
}

// NewTestServer creates a new TestServer in the given testing context and return its pointer
func NewTestServer(t *testing.T) *TestServer {
	return &TestServer{
		test:   t,
		router: NewRouter(NewTestErrorHandler(t)),
	}
}

// Init inits the TestServer's underlying httptest.Server with TestHandler's Router as its handler
func (ts *TestServer) Init() {
	ts.Server = NewServer(ts.router)
}

// Handler returns a new TestHandler initialized with the TestServer's testing context
func (ts *TestServer) Handler() *TestHandler {
	return NewTestHandler(NewTestErrorHandler(ts.test))
}

// Handle adds a TestHandler to the TestServer's Router
func (ts *TestServer) Handle(pathRegex, method string, handler http.Handler) {
	for _, route := range ts.router.routes {
		if route.regex == pathRegex {
			route.AddMethod(method, handler)
			return
		}
	}
	ts.router.AddRoute(NewRoute(pathRegex, NewTestErrorHandler(ts.test)).WithMethod(method, handler))
}
