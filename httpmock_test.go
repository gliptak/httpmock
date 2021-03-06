package httpmock

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/go/src/pkg/fmt"
	"github.com/stretchr/testify/require"
	"net/http"
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
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%v", res.StatusCode))
	}
	fb := new(FooBar)
	err = json.NewDecoder(res.Body).Decode(fb)
	return fb, err
}

func TestSingleOK(t *testing.T) {
	mock := setup(t)
	defer teardown()
	mock.AppendStep(Step{
		CheckRequest: func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/foobar", r.RequestURI)
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		},
		ReturnResponse: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fb := FooBar{Foo: "foo", Bar: 2}
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(fb)
			w.Write(b.Bytes())
		},
	})
	res, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Nil(t, err)
	require.Equal(t, "foo", res.Foo)
	require.Equal(t, 2, res.Bar)
}

func TestSingleForbidden(t *testing.T) {
	mock := setup(t)
	defer teardown()
	mock.AppendStep(Step{
		CheckRequest: func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/foobar", r.RequestURI)
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		},
		ReturnResponse: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		},
	})
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.EqualError(t, err, "403")
}

func TestSingleWrongFormat(t *testing.T) {
	mock := setup(t)
	defer teardown()
	mock.AppendStep(Step{
		CheckRequest: func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/foobar", r.RequestURI)
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		},
		ReturnResponse: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "<>")
		},
	})
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Contains(t, err.Error(), "invalid character ")
}

func TestSingleWildcard(t *testing.T) {
	mock := setup(t)
	defer teardown()
	mock.AppendStep(Step{
		CheckRequest: func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/foobar", r.RequestURI)
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		},
		ReturnResponse: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fb := FooBar{Foo: "foo", Bar: 2}
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(fb)
			w.Write(b.Bytes())
		},
	})
	res, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Nil(t, err)
	require.Equal(t, "foo", res.Foo)
	require.Equal(t, 2, res.Bar)
}

func TestEmptyStep(t *testing.T) {
	mock := setup(t)
	defer teardown()
	mock.AppendStep(Step{})
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.EqualError(t, err, "EOF")
}

func TestOutOfBounds1(t *testing.T) {
	setup(t)
	defer teardown()
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Contains(t, err, "EOF")
}

func TestOutOfBounds2(t *testing.T) {
	mock := setup(t)
	defer teardown()
	mock.AppendStep(Step{})
	_, _ = single(fmt.Sprintf("%s/foobar", server.URL))
	_, err := single(fmt.Sprintf("%s/foobar", server.URL))
	require.Contains(t, err, "EOF")
}
