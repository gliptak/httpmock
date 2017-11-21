package httpmock

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpMock struct {
	Handler []http.HandlerFunc
	current int
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) interface{}

func (mock *httpMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(){mock.current++}()
	mock.Handler[mock.current](w, r)
}

func NewHTTPMock() *httpMock { return new(httpMock) }

func (mock *httpMock) SetHandlers(handler []http.HandlerFunc) {
	mock.Handler = handler
}

var (
	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup(t *testing.T) *httpMock {
	// test server
	mock := NewHTTPMock()
	server = httptest.NewServer(mock)
	return mock
}

func teardown() {
	server.Close()
}
