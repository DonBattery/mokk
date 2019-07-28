package server

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// TestHandler is a pluggable struct capable of testing request headers and body
// and return predefined status, headers and body. If any error occurs during the handling of a request
// the TestHandler's ErrorHandler will be called with that error
type TestHandler struct {
	// Required properties
	requestHeaders http.Header
	requestBody    []byte

	// Response properties
	responseStatus  int
	responseHeaders http.Header
	responseBody    []byte

	errorHandler ErrorHandler
}

// NewTestHandler will create a new TestHandler with the given ErrorHandler
// if no ErrorHandler supplied it will fall back to the BasicErrorHandler
func NewTestHandler(errHandler ErrorHandler) *TestHandler {
	var handler ErrorHandler = &BasicErrorHandler{}
	if errHandler != nil {
		handler = errHandler
	}
	return &TestHandler{
		requestHeaders:  make(http.Header),
		responseHeaders: make(http.Header),
		errorHandler:    handler,
	}
}

// ServeHTTP
// The test handler will check (in order) for required request headers and body
// and call the ErrorHandler if and of them mismatches.
// Then it will write the response headers, status and body
func (handler *TestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check request headers
	if !containsAll(handler.requestHeaders, req.Header) {
		handler.errorHandler.HandleError(
			res,
			req,
			http.StatusBadRequest,
			errors.Errorf(
				"Required headers does not match with the actual headers.\nRequired:\n%+v\nActual:\n%+v\n",
				handler.requestHeaders,
				req.Header))
		return
	}
	// Check request body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.errorHandler.HandleError(
			res,
			req,
			http.StatusInternalServerError,
			errors.Wrap(err, "Cannot read request body"),
		)
		return
	}
	if !bytes.Equal(reqBody, handler.requestBody) {
		handler.errorHandler.HandleError(
			res,
			req,
			http.StatusBadRequest,
			errors.Errorf(
				"Required request body does not match with the actual request body.\nRequired:\n%s\nActual:\n%s\n",
				handler.requestBody,
				reqBody))
		return
	}
	// Write response headers
	if len(handler.responseHeaders) > 0 {
		for key, value := range handler.responseHeaders {
			for _, subValue := range value {
				res.Header().Add(key, subValue)
			}
		}
	}
	// Write response status code
	status := http.StatusOK
	if handler.responseStatus != 0 {
		status = handler.responseStatus
	}
	res.WriteHeader(status)
	// Write response body
	if handler.responseBody != nil {
		if _, err := res.Write(handler.responseBody); err != nil {
			handler.errorHandler.HandleError(
				res,
				req,
				http.StatusInternalServerError,
				errors.Wrap(err, "Failed to write response body"),
			)
			return
		}
	}
}

// WithRequestHeader adds a required request header to the TestHandler and returns it
func (handler *TestHandler) WithRequestHeader(key string, value string) *TestHandler {
	handler.requestHeaders.Add(key, value)
	return handler
}

// AddRequestHeader adds a required request header to the TestHandler
func (handler *TestHandler) AddRequestHeader(key string, value string) {
	handler.requestHeaders.Add(key, value)
}

// WithRequestHeaders adds all headers to the TestHandler and returns it
func (handler *TestHandler) WithRequestHeaders(headers map[string][]string) *TestHandler {
	for key, value := range headers {
		for _, subValue := range value {
			handler.requestHeaders.Add(key, subValue)
		}
	}
	return handler
}

// AddRequestHeaders adds all headers to the TestHandler
func (handler *TestHandler) AddRequestHeaders(headers map[string][]string) {
	for key, value := range headers {
		for _, subValue := range value {
			handler.requestHeaders.Add(key, subValue)
		}
	}
}

// WithRequestBody adds a required request body to the TestHandler and returns it
func (handler *TestHandler) WithRequestBody(body []byte) *TestHandler {
	handler.requestBody = body
	return handler
}

// AddRequestBody adds a required request body to the TestHandler
func (handler *TestHandler) AddRequestBody(body []byte) {
	handler.requestBody = body
}

// WithResponseBody adds a response body to the TestHandler and returns it
func (handler *TestHandler) WithResponseBody(body []byte) *TestHandler {
	handler.responseBody = body
	return handler
}

// AddResponseBody adds a response body to the TestHandler
func (handler *TestHandler) AddResponseBody(body []byte) {
	handler.responseBody = body
}

// AddResponseStatus adds a HTTP response status code to the TestHandler and returns it
func (handler *TestHandler) WithResponseStatus(statusCode int) *TestHandler {
	handler.responseStatus = statusCode
	return handler
}

// AddResponseStatus adds a HTTP response status code to the TestHandler
func (handler *TestHandler) AddResponseStatus(statusCode int) {
	handler.responseStatus = statusCode
}

// WithResponseHeader adds a response header to the TestHandler and returns it
func (handler *TestHandler) WithResponseHeader(key string, value string) *TestHandler {
	handler.responseHeaders.Add(key, value)
	return handler
}

// AddResponseHeader adds a response header to the TestHandler
func (handler *TestHandler) AddResponseHeader(key string, value string) {
	handler.responseHeaders.Add(key, value)
}

// WithResponseHeaders adds all response headers to the TestHandler and returns it
func (handler *TestHandler) WithResponseHeaders(headers map[string][]string) *TestHandler {
	for key, value := range headers {
		for _, subValue := range value {
			handler.responseHeaders.Add(key, subValue)
		}
	}
	return handler
}

// WithResponseHeaders adds all response headers to the TestHandler
func (handler *TestHandler) AddResponseHeaders(headers map[string][]string) {
	for key, value := range headers {
		for _, subValue := range value {
			handler.responseHeaders.Add(key, subValue)
		}
	}
}

// containsAll is a helper method which can tell if the actual request headers
// contains all the required request headers
func containsAll(required, actual map[string][]string) bool {
	if (required == nil) != (actual == nil) {
		return false
	}
	for requiredKey, requiredValue := range required {
		if _, ok := actual[requiredKey]; !ok {
			return false
		}
		for _, requiredSubValue := range requiredValue {
			if !stringInSlice(requiredSubValue, actual[requiredKey]) {
				return false
			}
		}
	}
	return true
}

// stringInSlice is a helper method to check if a string is in a slice of strings
func stringInSlice(required string, list []string) bool {
	for _, elem := range list {
		if elem == required {
			return true
		}
	}
	return false
}
