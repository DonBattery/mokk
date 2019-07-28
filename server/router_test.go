package server

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiMethod(t *testing.T) {
	t.Log("Testing multiple methods")

	srv := NewTestServer(t)
	srv.Handle("/test/route$", "GET", srv.Handler())
	srv.Handle("/test/route$", "POST", srv.Handler().WithResponseStatus(http.StatusCreated))
	srv.Init()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/test/route")
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response should be HTTP 200 OK on GET request")
	defer resp.Body.Close()

	postResp, postErr := http.Post(srv.URL+"/test/route", "application/json", nil)
	require.NoError(t, postErr, "The server souldn't return any aerrors")
	require.Equal(t, http.StatusCreated, postResp.StatusCode, "Response should be HTTP 201 Created on POST request")
	defer postResp.Body.Close()
}

func TestAddingRoutes(t *testing.T) {
	errHandler := NewTestErrorHandler(t)
	handler := NewTestHandler(errHandler)
	route := NewRoute("^/test$", errHandler)
	route.AddMethod("GET", handler)
	router := NewRouter(errHandler)
	router.AddRoute(route)
	srv := NewServer(router)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/test")
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response should be HTTP 200 OK on GET request")
	defer resp.Body.Close()
}

func TestOrder(t *testing.T) {
	t.Log("Testing order of routes")

	srv := NewTestServer(t)
	// This will only match for host/test/api and host/test/api/
	srv.Handle("^/test/api/?$", "GET", srv.Handler().WithResponseHeader("Route", "/test/api"))

	// This will match anything that ends with /api/users/123abc
	srv.Handle("/api/users/123abc$", "GET", srv.Handler().
		WithResponseHeader("Route", "/api/users/123abc").
		WithResponseBody([]byte("User found")))

	// This will match on any case
	srv.Handle(".*", "GET", srv.Handler().
		WithResponseStatus(http.StatusNotFound).
		WithResponseHeader("Route", "Not Found"))

	srv.Init()
	defer srv.Close()

	resp1, err1 := http.Get(srv.URL + "/test/api/users/123abc")
	require.NoError(t, err1, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp1.StatusCode, "Response should be HTTP 200 OK when getting existing user")
	defer resp1.Body.Close()
	respBody, readErr := ioutil.ReadAll(resp1.Body)
	require.NoError(t, readErr, "The response body should be readed")
	require.Equal(t, []byte("User found"), respBody, "The response body should be: \"User found\"")

	resp2, err2 := http.Get(srv.URL + "/test/api/")
	require.NoError(t, err2, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp2.StatusCode, "Response should be HTTP 200 OK on GET /test/api/")
	require.Equal(t, "/test/api", resp2.Header.Get("Route"), "Response should have the header Route:/test/api")
	defer resp2.Body.Close()

	resp3, err3 := http.Get(srv.URL + "/no/existent/route")
	require.NoError(t, err3, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusNotFound, resp3.StatusCode, "Response should be HTTP 404 on GET /test/api/")
	require.Equal(t, "Not Found", resp3.Header.Get("Route"), "Response should have the header Route:/test/api")
	defer resp3.Body.Close()
}
