package httpmock

import (
	"net/http"
	"net/http/httptest"
)

type HTTPMock struct {
	Handler []http.HandlerFunc
	current int
}

func (mock *HTTPMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(){mock.current++}()
	mock.Handler[mock.current](w, r)
}

func NewHTTPMock() *HTTPMock { return new(HTTPMock) }

func (mock *HTTPMock) SetHandlers(handler []http.HandlerFunc) {
	mock.Handler = handler
}

var (
	mock *HTTPMock

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup() {
	// test server
	mock = NewHTTPMock()
	server = httptest.NewServer(mock)
}

func teardown() {
	server.Close()
}
