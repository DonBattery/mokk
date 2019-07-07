package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Log("Testing simple GET request...")
	srv := NewTestServer(t)
	srv.Handle("/test/get$", "GET", srv.Handler())
	srv.Init()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/test/get")
	require.NoError(t, err, "Test server shouldn't return any errors")
	defer resp.Body.Close()
}

func TestMultipleGet(t *testing.T) {
	t.Log("Testing simple GET request...")
	srv := NewTestServer(t)
	srv.Handle("/test/get$", "GET", srv.Handler())
	srv.Handle("/test/get2$", "GET", srv.Handler())
	srv.Init()
	defer srv.Close()

	resp1, err := http.Get(srv.URL + "/test/get")
	require.NoError(t, err, "Test server shouldn't return any errors")
	defer resp1.Body.Close()

	resp2, err := http.Get(srv.URL + "/test/get2")
	require.NoError(t, err, "Test server shouldn't return any errors")
	defer resp2.Body.Close()
}

func TestStatus(t *testing.T) {
	t.Log("Testing return status...")
	srv := NewTestServer(t)
	srv.Handle("/test/status$", "GET", srv.Handler().WithResponseStatus(http.StatusCreated))
	srv.Init()
	defer srv.Close()

	res, err := http.Get(srv.URL + "/test/status")
	require.NoError(t, err, "Test server shouldn't return any errors")
	defer res.Body.Close()
	require.Equal(t, http.StatusCreated, res.StatusCode, "Response status should be 201 Created")
}

func TestHeader(t *testing.T) {
	t.Log("Testing GET request with required header...")
	srv := NewTestServer(t)
	srv.Handle("/test/header$", "GET", srv.Handler().WithRequestHeader("Auth", "Pass"))
	srv.Init()
	defer srv.Close()

	client := http.Client{}
	request, err := http.NewRequest("GET", srv.URL+"/test/header", nil)
	require.NoError(t, err, "Test request should be created")
	request.Header.Add("Auth", "Pass")
	resp, getErr := client.Do(request)
	require.NoError(t, getErr, "Test server shouldn't return any errors")
	defer resp.Body.Close()
}
