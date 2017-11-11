package httpmock

import (
	"bytes"
	"encoding/json"
	"github.com/golang/go/src/pkg/fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"github.com/stretchr/testify/require"
)

type FooBar struct {
	Foo string
	Bar int
}

func single(url string) (*FooBar, error) {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%v", res.StatusCode))
	}
	fb := new(FooBar)
	err = json.NewDecoder(res.Body).Decode(fb)
	return fb, err
}

func TestSingleOK(t *testing.T) {
	setup()
	defer teardown()
	mock.SetHandler(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foobar", r.RequestURI)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fb := FooBar{Foo: "foo", Bar: 2}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(fb)
		w.Write(b.Bytes())
	})
	res, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Nil(t, err)
	require.Equal(t, "foo", res.Foo)
	require.Equal(t, 2, res.Bar)
}

func TestSingleForbidden(t *testing.T) {
	setup()
	defer teardown()
	mock.SetHandler(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foobar", r.RequestURI)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusForbidden)
	})
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.EqualError(t, err, "403")
}

func TestSingleWrongFormat(t *testing.T) {
	setup()
	defer teardown()
	mock.SetHandler(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foobar", r.RequestURI)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "<>")
	})
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Contains(t, err.Error(), "invalid character ")
}

func TestSingleWildcard(t *testing.T) {
	setup()
	defer teardown()
	mock.SetHandler(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foobar", r.RequestURI)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fb := FooBar{Foo: "foo", Bar: 2}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(fb)
		w.Write(b.Bytes())
	})
	res, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Nil(t, err)
	require.Equal(t, "foo", res.Foo)
	require.Equal(t, 2, res.Bar)
}

type HTTPMock struct {
	Handler http.HandlerFunc
}

func (mock *HTTPMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mock.Handler(w, r)
}

func NewHTTPMock() *HTTPMock { return new(HTTPMock) }

func (mock *HTTPMock) SetHandler(handler func(http.ResponseWriter, *http.Request)) {
	mock.Handler = http.HandlerFunc(handler)
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
