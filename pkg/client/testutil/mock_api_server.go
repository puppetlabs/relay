package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

type mockRoute struct {
	path       string
	response   interface{}
	statusCode int
	header     map[string]string
}

type MockRoutes struct {
	routes map[string]*mockRoute
}

func (mr *MockRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := mr.routes[r.URL.Path]

	if !ok {
		http.NotFound(w, r)

		return
	}

	b, err := json.Marshal(route.response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	buf := bytes.NewReader(b)

	w.WriteHeader(route.statusCode)
	w.Header().Set("content-type", "application/json")

	if route.header != nil {
		for k, v := range route.header {
			w.Header().Add(k, v)
		}
	}

	io.Copy(w, buf)
}

func (mr *MockRoutes) Add(path string, status int, resp interface{}, header map[string]string) {
	if mr.routes == nil {
		mr.routes = make(map[string]*mockRoute)
	}

	mr.routes[path] = &mockRoute{
		path:       path,
		response:   resp,
		statusCode: status,
		header:     header,
	}
}

func WithTestServer(h http.Handler, fn func(ts *httptest.Server)) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// wrapped in HandlerFunc for debugging
		h.ServeHTTP(w, r)
	}))

	defer ts.Close()

	fn(ts)
}
