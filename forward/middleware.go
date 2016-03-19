package forward

import (
	"net/http"
	"net/url"
)

// TODO: support passing options
func To(uri string) func(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	fwd, err := New(PassHostHeader(true))
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = parsedURL.Scheme
		r.URL.Host = parsedURL.Host
		r.Host = parsedURL.Host

		// Forward the HTTP request
		fwd.ServeHTTP(w, r)
	}
}
