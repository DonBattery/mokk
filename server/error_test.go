package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicErrorHandler_NoError(t *testing.T) {
	handler := NewTestHandler(nil)
	route := NewRoute("/test/error$", nil)
	route.AddMethod("GET", handler)
	router := NewRouter(nil)
	router.AddRoute(route)
	server := NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test/error")
	require.NoError(t, err, "GET shouldn't return any errors")
	defer resp.Body.Close()
}

func TestBasicErrorHandler_NotFound(t *testing.T) {
	server := NewServer(NewRouter(nil))
	defer server.Close()

	resp, err := http.Get(server.URL + "/non/existent/path")
	require.NoError(t, err, "GET shouldn't return any errors")
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Status should be 404 Not Found")
}

func TestBasicErrorHandler_MethodNotAllowed(t *testing.T) {
	handler := NewTestHandler(nil)
	route := NewRoute("/existing/path$", nil)
	route.AddMethod("PUT", handler)
	router := NewRouter(nil)
	router.AddRoute(route)
	server := NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/existing/path")
	require.NoError(t, err, "GET shouldn't return any errors")
	defer resp.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Status should be 405 Method Not Allowed")
}
