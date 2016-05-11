package vinxi

import (
	"bytes"
	"fmt"
	"github.com/nbio/st"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVinxi(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	}))
	defer ts.Close()

	v := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", ts.URL, bytes.NewBufferString("foo"))

	v.ServeHTTP(w, req)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Body.String(), "Hello world\n")
}
