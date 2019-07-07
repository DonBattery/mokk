package server

import (
	"fmt"
	"net/http"
	"testing"
)

type ErrorHandler interface {
	HandleError(res http.ResponseWriter, req *http.Request, status int, err error)
}

type BasicErrorHandler struct{}

func (h *BasicErrorHandler) HandleError(res http.ResponseWriter, req *http.Request, status int, err error) {
	res.WriteHeader(status)
	if err != nil {
		_, writeErr := res.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
		panic(writeErr)
	}
}

type TestErrorHandler struct {
	T *testing.T
}

func NewTestErrorHandler(t *testing.T) *TestErrorHandler {
	return &TestErrorHandler{
		T: t,
	}
}

func (h *TestErrorHandler) HandleError(res http.ResponseWriter, req *http.Request, status int, err error) {
	h.T.Errorf("HTTP response status: %d Error: %s", status, err)
}
