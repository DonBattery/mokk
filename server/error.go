package server

import (
	"fmt"
	"net/http"
	"testing"
)

// ErrorHandler is a struct that which can HandleError s
// It extends the handler arguments (the response writer and the request)
// with a HTTP status code and an error
type ErrorHandler interface {
	HandleError(res http.ResponseWriter, req *http.Request, status int, err error)
}

// BasicErrorHandler is the default error handler
type BasicErrorHandler struct{}

// HandleError writes the status code on the response's header and the error string in the response's body
func (h *BasicErrorHandler) HandleError(res http.ResponseWriter, req *http.Request, status int, err error) {
	res.WriteHeader(status)
	if err != nil {
		_, writeErr := fmt.Fprintf(res, "Status: %d Error: %s\n", status, err.Error())
		if writeErr != nil {
			panic(writeErr)
		}
	}
}

// TestErrorHandler is the default error handler of the TestServer
type TestErrorHandler struct {
	T *testing.T
}

// NewTestErrorHandler creates a TestErrorHandler in the given testing context
// and returns its pointer
func NewTestErrorHandler(t *testing.T) *TestErrorHandler {
	return &TestErrorHandler{
		T: t,
	}
}

// HandleError will Logf the error message than Fail the test
func (h *TestErrorHandler) HandleError(res http.ResponseWriter, req *http.Request, status int, err error) {
	h.T.Errorf("HTTP response status: %d Error: %s", status, err)
}
