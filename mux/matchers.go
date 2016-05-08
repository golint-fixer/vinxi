package mux

import (
	"net/http"
	"regexp"
)

// Matcher represents the function interface implemented by matchers
type Matcher func(req *http.Request) bool

// MatchMethod matches the HTTP method name againts the request.
func MatchMethod(methods ...string) Matcher {
	return func(req *http.Request) bool {
		for _, method := range methods {
			if req.Method == method {
				return true
			}
		}
		return false
	}
}

// MatchPath matches the given path patterns againts the incoming request.
func MatchPath(pattern string) Matcher {
	rex := regexp.MustCompile(pattern)
	return func(req *http.Request) bool {
		path := req.URL.Path
		if pattern == path {
			return true
		}
		return rex.MatchString(path)
	}
}

// MatchHost matches the given host string in the incoming request.
func MatchHost(pattern string) Matcher {
	rex := regexp.MustCompile(pattern)
	return func(req *http.Request) bool {
		if req.URL.Host == pattern {
			return true
		}
		return rex.MatchString(req.Host)
	}
}

// MatchQuery matches a given query param againts the request.
func MatchQuery(key, pattern string) Matcher {
	rex := regexp.MustCompile(pattern)
	return func(req *http.Request) bool {
		return rex.MatchString(req.URL.Query().Get(key))
	}
}

// MatchHeader matches a given header key and value againts the request.
func MatchHeader(key, pattern string) Matcher {
	rex := regexp.MustCompile(pattern)
	return func(req *http.Request) bool {
		return rex.MatchString(req.Header.Get(key))
	}
}
