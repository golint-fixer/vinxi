package vinxi

import (
	"fmt"
	"gopkg.in/h2non/baloo.v0"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVinxi(t *testing.T) {
	v := New().Use(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	}))
	ts := httptest.NewServer(v)
	baloo.New(ts.URL).Get("/").Expect(t).StatusOk().BodyEquals("Hello world").Done()
}

func TestRouteGET(t *testing.T) {
	v := New()
	v.Get("/hello").Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	}))
	ts := httptest.NewServer(v)
	baloo.New(ts.URL).Get("/hello").Expect(t).StatusOk().BodyEquals("Hello world").Done()
	baloo.New(ts.URL).Get("/").Expect(t).Status(502).BodyEquals("Bad Gateway").Done()
}
