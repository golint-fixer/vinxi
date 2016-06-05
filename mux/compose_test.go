package mux

import (
	"github.com/nbio/st"
	"gopkg.in/vinxi/vinxi.v0/utils"
	"net/http"
	"net/url"
	"testing"
)

type composer func(...*Mux) *Mux

func testMuxIfMatches(t *testing.T, c composer) {
	mx := New()
	mx.Use(c(Method("GET"), Host("foo.com")).Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	}))

	wrt := utils.NewWriterStub()
	req := newRequest()
	req.URL.Host = "foo.com"

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}

func TestMuxComposeIfMatches(t *testing.T) {
	testMuxIfMatches(t, If)
}

func TestMuxComposeEveryMatches(t *testing.T) {
	testMuxIfMatches(t, Every)
}

func testMuxComposeIfUnmatch(t *testing.T, c composer) {
	mx := New()
	mx.Use(c(Method("GET"), Host("bar.com")).Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	}))

	wrt := utils.NewWriterStub()
	req := newRequest()
	req.URL.Host = "foo.com"

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "")
}

func TestMuxComposeIfUnmatch(t *testing.T) {
	testMuxComposeIfUnmatch(t, If)
}

func TestMuxComposeEveryUnmatch(t *testing.T) {
	testMuxComposeIfUnmatch(t, Every)
}

func testMuxComposeOrMatch(t *testing.T, c composer) {
	mx := New()
	mx.Use(c(Method("GET"), Host("bar.com")).Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	}))

	wrt := utils.NewWriterStub()
	req := newRequest()
	req.URL.Host = "foo.com"

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}

func TestMuxComposeOrMatch(t *testing.T) {
	testMuxComposeOrMatch(t, Or)
}

func TestMuxComposeSomeMatch(t *testing.T) {
	testMuxComposeOrMatch(t, Some)
}

func testMuxComposeOrUnmatch(t *testing.T, c composer) {
	mx := New()
	mx.Use(c(Method("GET"), Host("bar.com")).Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	}))

	wrt := utils.NewWriterStub()
	req := newRequest()
	req.URL.Host = "foo.com"

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")

	req.Method = "POST"
	wrt = utils.NewWriterStub()
	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "")
}

func TestMuxComposeOrUnmatch(t *testing.T) {
	testMuxComposeOrUnmatch(t, Or)
}

func TestMuxComposeSomeUnmatch(t *testing.T) {
	testMuxComposeOrUnmatch(t, Some)
}

func newRequest() *http.Request {
	return &http.Request{URL: &url.URL{}, Header: make(http.Header), Method: "GET"}
}
