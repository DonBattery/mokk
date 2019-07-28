package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainsAll(t *testing.T) {
	t.Log("Testing containsAll...")

	required := map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
	}
	actual := map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
			"Elem3",
		},
		"List2": {
			"Elem1",
			"Elem2",
		},
	}
	require.True(
		t,
		containsAll(required, actual),
		"containsAll should return true if required elems are in the actual map")

	actual = map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
	}
	required = map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
		"List2": {
			"Elem1",
			"Elem2",
		},
	}
	require.False(
		t,
		containsAll(required, actual),
		"containsAll should return false if not all required elem are in the actual map")

	actual = map[string][]string{
		"List1": {
			"Elem1",
		},
	}
	required = map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
	}
	require.False(
		t,
		containsAll(required, actual),
		"containsAll should return false if not all required elem are in the actual map")

	require.True(t, containsAll(nil, nil), "containsAll should return true if both maps are nil")

	require.False(t, containsAll(map[string][]string{
		"Key": {
			"Value1", "Value2",
		},
	}, nil), "containsAll should return false if the required map is declared but the actual is nil")

	require.False(t, containsAll(
		nil,
		map[string][]string{
			"Key": {
				"Value1", "Value2",
			}}),
		"containsAll should return false if the required map is nil but the actual map is declared")
}

func TestAddResponseStatus(t *testing.T) {
	t.Log("Testing AddResponseStatus..")

	handler := NewTestHandler(nil)
	handler.AddResponseStatus(http.StatusTeapot)
	srv := NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusTeapot, resp.StatusCode, "Response status should be HTTP 418 I'm a Teapot")
	defer resp.Body.Close()
}

func TestRequestHeaders_bad_request(t *testing.T) {
	t.Log("Testing TestHandler request headers on bad request...")

	handler := NewTestHandler(nil)
	handler.AddRequestHeader("Key", "Value")
	srv := NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response should be HTTP 400 Bad Request")
	defer resp.Body.Close()
}

func TestRequestHeaders_multiple_headers(t *testing.T) {
	t.Log("Testing TestHandler request headers on multiple headers...")

	handler := NewTestHandler(nil).WithRequestHeaders(map[string][]string{
		"Key1": {"k1-Value1", "k1-Value2"},
		"Key2": {"k2-Value1", "k2-Value2"},
	})
	handler.AddRequestHeader("Key1", "k1-Value3")
	handler.AddRequestHeader("Key3", "k3-Value1")
	handler.AddRequestHeaders(map[string][]string{
		"Key3": {"k3-Value2", "k3-Value3"},
		"Key4": {"k4-Value1"},
	})
	srv := NewServer(handler)
	defer srv.Close()

	client := http.Client{}
	request, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err, "Test request should be created")
	request.Header.Add("Key1", "k1-Value1")
	request.Header.Add("Key1", "k1-Value2")
	request.Header.Add("Key1", "k1-Value3")
	request.Header.Add("Key2", "k2-Value1")
	request.Header.Add("Key2", "k2-Value2")
	request.Header.Add("Key3", "k3-Value1")
	request.Header.Add("Key3", "k3-Value2")
	request.Header.Add("Key3", "k3-Value3")
	request.Header.Add("Key4", "k4-Value1")
	resp, getErr := client.Do(request)
	require.NoError(t, getErr, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response should be HTTP 200 OK")
	defer resp.Body.Close()
}

func TestResponseHeader(t *testing.T) {
	t.Log("Testing response headers...")

	handler := NewTestHandler(nil)
	handler.AddResponseHeader("Key1", "Value1")
	srv := NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be HTTP 200 OK")
	require.Equal(t, "Value1", resp.Header.Get("Key1"), "The respons should have the header Key1:Value1")
	defer resp.Body.Close()
}

func TestResponseHeaders(t *testing.T) {
	t.Log("Testing response headers...")

	handler := NewTestHandler(nil).
		WithResponseHeader("Key1", "k1-Value1").
		WithResponseHeader("Key1", "k1-Value2").
		WithResponseHeaders(map[string][]string{
			"Key1": {"k1-Value3", "k1-Value4"},
			"Key2": {"k2-Value1", "k2-Value2"},
		})
	handler.AddResponseHeader("Key1", "k1-Value5")
	handler.AddResponseHeader("Key3", "k3-Value1")
	handler.AddResponseHeaders(map[string][]string{
		"Key1": {"k1-Value6", "k1-Value7"},
		"Key4": {"k4-Value1"},
	})
	srv := NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be HTTP 200 OK")
	requiredHeaders := http.Header{
		"Key1": {"k1-Value1", "k1-Value2", "k1-Value3", "k1-Value4", "k1-Value5", "k1-Value6", "k1-Value7"},
		"Key2": {"k2-Value1", "k2-Value2"},
		"Key3": {"k3-Value1"},
		"Key4": {"k4-Value1"},
	}
	require.Truef(t,
		containsAll(requiredHeaders, resp.Header),
		"The respons should contain all the supplyed headers\nRequired:\n%+v\nActual:\n%+v\n",
		requiredHeaders,
		resp.Header,
	)

	defer resp.Body.Close()

}

func TestAddRequestBody(t *testing.T) {
	t.Log("Testing AddReqestBody...")

	handler := NewTestHandler(nil)
	handler.AddRequestBody([]byte("Test Body"))
	srv := NewServer(handler)
	defer srv.Close()

	client := http.Client{}
	request, err := http.NewRequest("GET", srv.URL, bytes.NewReader([]byte("Test Body")))
	require.NoError(t, err, "The request should be created")
	resp, getErr := client.Do(request)
	require.NoError(t, getErr, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response should be HTTP 200 OK")
	defer resp.Body.Close()
}

func TestRequestBody_bad_request(t *testing.T) {
	t.Log("Testing missing request body...")

	handler := NewTestHandler(nil)
	handler.AddRequestBody([]byte("Test Body"))
	srv := NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response status should be HTTP 400 Bad Request")

	respBody, readErr := ioutil.ReadAll(resp.Body)
	require.NoError(t, readErr, "The respons body should be readable")
	requiredError := "Required request body does not match with the actual request body"
	require.Truef(t,
		strings.Contains(string(respBody), requiredError),
		"Expected response should contain:\n%s\nActual Response:\n%s\n", requiredError, respBody)
	defer resp.Body.Close()
}

func TestResponseBody(t *testing.T) {
	t.Log("Testing response body...")

	handler := NewTestHandler(nil)
	handler.AddResponseBody([]byte("Test Body"))
	srv := NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err, "Test server shouldn't return any errors")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be HTTP 200 OK")
	respBody, readErr := ioutil.ReadAll(resp.Body)
	require.NoError(t, readErr, "The response body should be readable")
	require.Equal(t, []byte("Test Body"), respBody, "The response body should be: \"Test Body\"")
	defer resp.Body.Close()
}
