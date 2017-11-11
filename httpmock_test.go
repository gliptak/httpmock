package httpmock

import (
	"bytes"
	"encoding/json"
	"github.com/golang/go/src/pkg/fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
	fb := new(FooBar)
	err = json.NewDecoder(res.Body).Decode(fb)
	return fb, err
}

func TestSingle(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/foobar", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fb := FooBar{Foo: "foo", Bar: 2}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(fb)
		w.Write(b.Bytes())
	})
	res, err := single(fmt.Sprintf("%s/foobar", server.URL))
	assert.Nil(t, err)
	assert.Equal(t, "foo", res.Foo)
	assert.Equal(t, 2, res.Bar)
}

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func teardown() {
	server.Close()
}
