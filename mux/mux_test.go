package mux

import (
	"github.com/nbio/st"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	fooMw = func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	}
	barMw = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("bar", "foo")
	})
)

func TestMux_Match_matched(t *testing.T) {
	mx := New()
	mx.AddMatcher(MatchMethod("GET"))

	req := newRequest()
	req.Method = "GET"

	st.Expect(t, mx.Match(req), true)
}

func TestMux_Match_unmatched(t *testing.T) {
	mx := New()
	mx.AddMatcher(MatchMethod("POST"))

	req := newRequest()
	req.Method = "GET"

	st.Expect(t, mx.Match(req), false)
}

func testMuxAddMatcherAndIf(t *testing.T, f func(m *Mux) func(matchers ...Matcher) *Mux) {
	mx := New()
	f(mx)(MatchMethod("POST"), MatchHost("www.example.com"))

	req := newRequest()
	req.Method = "POST"
	st.Expect(t, mx.Match(req), false)

	req.Host = "www.example.com"
	st.Expect(t, mx.Match(req), true)
}

func TestMux_Addmatcher(t *testing.T) {
	testMuxAddMatcherAndIf(t, func(m *Mux) func(matchers ...Matcher) *Mux {
		return m.AddMatcher
	})
}

func TestMux_If(t *testing.T) {
	testMuxAddMatcherAndIf(t, func(m *Mux) func(matchers ...Matcher) *Mux {
		return m.If
	})
}

func TestMux_Some(t *testing.T) {
	mx := New()
	mx.Some(MatchMethod("POST"), MatchHost("www.example.com"))

	req := newRequest()

	st.Expect(t, mx.Match(req), false)

	req.Method = "POST"
	st.Expect(t, mx.Match(req), true)

	req.Method = "GET"
	req.Host = "www.example.com"
	st.Expect(t, mx.Match(req), true)

	req.Method = "GET"
	st.Expect(t, mx.Match(req), true)
}

func TestMux_Use(t *testing.T) {
	mx := New()
	mx.Use(fooMw)

	wrt := httptest.NewRecorder()
	req := newRequest()

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}

func TestMux_UsePhase(t *testing.T) {
	mx := New()
	mx.UsePhase("request", fooMw)

	wrt := httptest.NewRecorder()
	req := newRequest()

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}

func TestMux_UseFinalHandler(t *testing.T) {
	mx := New()
	mx.Use(fooMw)
	mx.UseFinalHandler(barMw)

	wrt := httptest.NewRecorder()
	req := newRequest()

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
	st.Expect(t, wrt.Header().Get("bar"), "foo")
}

func TestMux_HandleHTTP_matched(t *testing.T) {
	mx := New()
	mx.If(MatchMethod("GET"))
	mx.Use(fooMw)

	wrt := httptest.NewRecorder()
	req := newRequest()
	req.Method = "GET"

	mx.HandleHTTP(wrt, req, barMw)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
	st.Expect(t, wrt.Header().Get("bar"), "foo")
}

func TestMux_HandleHTTP_unmatched(t *testing.T) {
	mx := New()
	mx.If(MatchMethod("GET"))
	mx.Use(fooMw)

	wrt := httptest.NewRecorder()
	req := newRequest()
	req.Method = "POST"

	mx.HandleHTTP(wrt, req, barMw)
	st.Expect(t, wrt.Header().Get("bar"), "foo")
}
