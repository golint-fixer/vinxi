package mux

// Match creates a new multiplexer based on a given matcher function.
func Match(matchers ...Matcher) *Mux {
	mx := New()
	mx.AddMatcher(matchers...)
	return mx
}

// Method returns a new multiplexer who matches an HTTP request based on the given method/s.
func Method(methods ...string) *Mux {
	return Match(MatchMethod(methods...))
}

// Path returns a new multiplexer who matches an HTTP request
// path based on the given regexp pattern.
func Path(pattern string) *Mux {
	return Match(MatchPath(pattern))
}

// Host returns a new multiplexer who matches an HTTP request
// URL host based on the given regexp pattern.
func Host(pattern string) *Mux {
	return Match(MatchHost(pattern))
}

// Query returns a new multiplexer who matches an HTTP request
// query param based on the given key and regexp pattern.
func Query(key, pattern string) *Mux {
	return Match(MatchQuery(key, pattern))
}

// Header returns a new multiplexer who matches an HTTP request
// header field based on the given key and regexp pattern.
func Header(key, pattern string) *Mux {
	return Match(MatchHeader(key, pattern))
}
