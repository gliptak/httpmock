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
		mock.t.Fatalf("step out of bounds %v <= %v\n", len(mock.steps), mock.current)
		//mock.t.Fail()
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
