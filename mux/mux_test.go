package mux

import (
	"github.com/nbio/st"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMuxSimple(t *testing.T) {
	mx := New()
	mx.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	})

	wrt := httptest.NewRecorder()
	req := newRequest()

	mx.Layer.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}
