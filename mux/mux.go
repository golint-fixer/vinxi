// Package mux implements an HTTP domain-specific traffic multiplexer
// with built-in matchers and features for easy plugin composition and activable logic.
package mux

import (
	"gopkg.in/vinxi/layer.v0"
	"net/http"
)

// Mux is a HTTP request/response/error multiplexer who implements both
// middleware and plugin interfaces.
// It has been designed for easy plugin composition based on HTTP matchers/filters.
type Mux struct {
	// Matchers stores a list of matcher functions.
	Matchers []Matcher

	// Layer stores the multiplexer middleware layer.
	Layer *layer.Layer
}

// New creates a new multiplexer with default settings.
func New() *Mux {
	return &Mux{Layer: layer.New()}
}

// Match matches the give Context againts a list of matchers and
// returns `true` if all the matchers passed.
func (m *Mux) Match(req *http.Request) bool {
	for _, matcher := range m.Matchers {
		if !matcher(req) {
			return false
		}
	}
	return true
}

// AddMatcher adds a new matcher function in the current mumultiplexer matchers stack.
func (m *Mux) AddMatcher(matchers ...Matcher) *Mux {
	m.Matchers = append(m.Matchers, matchers...)
	return m
}

// If is a semantic alias to AddMatcher.
func (m *Mux) If(matchers ...Matcher) *Mux {
	return m.AddMatcher(matchers...)
}

// Some matches the incoming request if at least one of the matchers passes.
func (m *Mux) Some(matchers ...Matcher) *Mux {
	return m.AddMatcher(func(req *http.Request) bool {
		for _, matcher := range matchers {
			if matcher(req) {
				return true
			}
		}
		return false
	})
}

// Use registers a new plugin in the middleware stack.
func (m *Mux) Use(handler interface{}) *Mux {
	m.Layer.Use(layer.RequestPhase, handler)
	return m
}

// UsePhase registers a new plugin in the middleware stack.
func (m *Mux) UsePhase(phase string, handler interface{}) *Mux {
	m.Layer.Use(phase, handler)
	return m
}

// UseFinalHandler registers a new plugin in the middleware stack.
func (m *Mux) UseFinalHandler(handler http.Handler) *Mux {
	m.Layer.UseFinalHandler(handler)
	return m
}

// HandleHTTP returns the function handler to match an incoming HTTP transacion
// and trigger the equivalent middleware phase.
func (m *Mux) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	if m.Match(r) {
		m.Layer.Run("request", w, r, h)
		return
	}
	h.ServeHTTP(w, r)
}
