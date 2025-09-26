package httpmock

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpMock struct {
	t       *testing.T
	steps   []Step
	current int
}

type Step struct {
	CheckRequest   func(w http.ResponseWriter, r *http.Request)
	ReturnResponse func(w http.ResponseWriter, r *http.Request)
}

func (mock *httpMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { mock.current++ }()
	if len(mock.steps) <= mock.current {
		// Instead of calling Fatalf (which kills the test), close the connection
		// This will cause the client to receive an EOF error
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, err := hj.Hijack()
			if err == nil {
				conn.Close()
				return
			}
		}
		// Fallback: write nothing and return, which should cause connection issues
		return
	}
	if mock.steps[mock.current].CheckRequest != nil {
		mock.steps[mock.current].CheckRequest(w, r)
	}
	if mock.steps[mock.current].ReturnResponse != nil {
		mock.steps[mock.current].ReturnResponse(w, r)
	}
}

func NewHTTPMock(t *testing.T) *httpMock {
	mock := new(httpMock)
	mock.t = t
	return mock
}

func (mock *httpMock) AppendStep(step Step) {
	mock.steps = append(mock.steps, step)
}

var (
	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup(t *testing.T) *httpMock {
	// test server
	mock := NewHTTPMock(t)
	server = httptest.NewServer(mock)
	return mock
}

func teardown() {
	server.Close()
}
