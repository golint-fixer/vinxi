package mux

import (
	"github.com/nbio/st"
	"gopkg.in/vinxi/utils.v0"
	"net/http"
	"net/url"
	"testing"
)

func TestMethod(t *testing.T) {
	mx := New()
	mx.Use(Method("GET").Use(pass))

	wrt := utils.NewWriterStub()
	req := newRequest()

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}

func TestPath(t *testing.T) {
	cases := []struct {
		value   string
		path    string
		matches bool
	}{
		{"baz", "/bar/foo/baz", true},
		{"bar", "/bar/foo/baz", true},
		{"^/bar", "/bar/foo/baz", true},
		{"/foo/", "/bar/foo/baz", true},
		{"f*", "/bar/foo/baz", true},
		{"fo[o]", "/bar/foo/baz", true},
		{"baz$", "/bar/foo/baz", true},
		{"foo$", "/bar/foo/baz", false},
		{"foobar", "/bar/foo/baz", false},
	}

	for _, test := range cases {
		mx := New()
		mx.Use(Path(test.value).Use(pass))
		wrt := utils.NewWriterStub()
		req := newRequest()
		req.URL.Path = test.path
		mx.Layer.Run("request", wrt, req, nil)
		match(t, wrt, test.matches)
	}
}

func TestHost(t *testing.T) {
	cases := []struct {
		value   string
		url     string
		matches bool
	}{
		{"foo.com", "http://foo.com", true},
		{"foo", "http://foo.com", true},
		{".com", "http://foo.com", true},
		{"foo.com$", "http://foo.com", true},
		{"bar", "http://foo.com", false},
		{"^http://foo", "http://foo.com", false},
		{"^foo.com$", "http://foo.com", true},
	}

	for _, test := range cases {
		mx := New()
		mx.Use(Host(test.value).Use(pass))
		wrt := utils.NewWriterStub()
		req := newRequest()
		req.URL, _ = url.Parse(test.url)
		req.Host = req.URL.Host
		mx.Layer.Run("request", wrt, req, nil)
		match(t, wrt, test.matches)
	}
}

func TestQuery(t *testing.T) {
	cases := []struct {
		key     string
		value   string
		url     string
		matches bool
	}{
		{"foo", "bar", "http://baz.com?foo=bar", true},
		{"foo", "^bar$", "http://foo.com?foo=bar&baz=foo", true},
		{"foo", "b[a]r", "http://foo.com?foo=bar&baz=foo", true},
		{"foo", "foo", "http://foo.com?foo=bar&baz=foo", false},
		{"foo", "baz", "http://foo.com?foo=bar&baz=foo", false},
		{"baz", "foo", "http://foo.com?foo=bar&baz=foo", true},
		{"foo", "foo", "http://foo.com", false},
	}

	for _, test := range cases {
		mx := New()
		mx.Use(Query(test.key, test.value).Use(pass))
		wrt := utils.NewWriterStub()
		req := newRequest()
		req.URL, _ = url.Parse(test.url)
		mx.Layer.Run("request", wrt, req, nil)
		match(t, wrt, test.matches)
	}
}

func TestHeader(t *testing.T) {
	cases := []struct {
		key     string
		value   string
		headers map[string]string
		matches bool
	}{
		{"foo", "bar", map[string]string{"foo": "bar"}, true},
		{"foo", "bar", map[string]string{"foo": "foobar"}, true},
		{"foo", "bar", map[string]string{"foo": "foo"}, false},
		{"foo", "bar", map[string]string{}, false},
	}

	for _, test := range cases {
		mx := New()
		mx.Use(Header(test.key, test.value).Use(pass))
		wrt := utils.NewWriterStub()
		req := newRequest()
		for key, value := range test.headers {
			req.Header.Set(key, value)
		}
		mx.Layer.Run("request", wrt, req, nil)
		match(t, wrt, test.matches)
	}
}

func pass(w http.ResponseWriter, r *http.Request, h http.Handler) {
	w.Header().Set("foo", "bar")
	h.ServeHTTP(w, r)
}

func match(t *testing.T, w http.ResponseWriter, shouldMatch bool) {
	if shouldMatch {
		st.Expect(t, w.Header().Get("foo"), "bar")
	} else {
		st.Expect(t, w.Header().Get("foo"), "")
	}
}
