package mux

import (
	"github.com/nbio/st"
	"net/url"
	"testing"
)

func TestMatchMethod(t *testing.T) {
	req := newRequest()
	st.Expect(t, MatchMethod("GET")(req), true)

	req = newRequest()
	req.Method = "POST"
	st.Expect(t, MatchMethod("GET")(req), false)

	req = newRequest()
	st.Expect(t, MatchMethod("POST", "GET")(req), true)
}

func TestMatchPath(t *testing.T) {
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
		req := newRequest()
		req.URL.Path = test.path
		st.Expect(t, MatchPath(test.value)(req), test.matches)
	}
}

func TestMatchHost(t *testing.T) {
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
		req := newRequest()
		req.URL, _ = url.Parse(test.url)
		req.Host = req.URL.Host
		st.Expect(t, MatchHost(test.value)(req), test.matches)
	}
}

func TestMatchQuery(t *testing.T) {
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
		req := newRequest()
		req.URL, _ = url.Parse(test.url)
		st.Expect(t, MatchQuery(test.key, test.value)(req), test.matches)
	}
}

func TestMatchRequestHeader(t *testing.T) {
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
		req := newRequest()
		for key, value := range test.headers {
			req.Header.Set(key, value)
		}
		st.Expect(t, MatchHeader(test.key, test.value)(req), test.matches)
	}
}
