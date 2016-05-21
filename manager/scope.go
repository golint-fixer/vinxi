package manager

import (
	"net/http"

	"github.com/dchest/uniuri"
)

// Rule represents the required interface implemented
// by HTTP traffic rules.
//
// Rule is designed to inspect an incoming HTTP
// traffic and determine if should trigger the registered
// plugins if the rule matches.
type Rule interface {
	// ID returns the rule unique identifier.
	ID() string
	// Name returns the rule semantic alias.
	Name() string
	// Description is used to retrieve the rule semantic description.
	Description() string
	// JSONConfig is used to retrieve the rule config as JSON.
	JSONConfig() string
	// Match is used to determine if a given http.Request
	// passes the rule assertion.
	Match(*http.Request) bool
}

// Scope represents the HTTP configuration scope who can
// store rules and plugins.
type Scope struct {
	disabled bool
	rules    []Rule
	plugins  *PluginLayer

	ID          string
	Name        string
	Description string
}

// NewScope creates a new Scope instance
// with the given name alias and optional description.
func NewScope(name, description string) *Scope {
	return &Scope{ID: uniuri.New(), Name: name, Description: description, plugins: NewPluginLayer()}
}

func (s *Scope) UsePlugin(plugins ...Plugin) {
	s.plugins.Use(plugins...)
}

func (s *Scope) UseRule(rules ...Rule) {
	s.rules = append(s.rules, rules...)
}

func (s *Scope) Rules() []Rule {
	return s.rules
}

func (s *Scope) RemoveRule(id string) bool {
	return true
}

func (s *Scope) Disable() {
	s.disabled = true
}

func (s *Scope) Enable() {
	s.disabled = false
}

func (s *Scope) IsEnabled() bool {
	return s.disabled == false
}

func (s *Scope) HandleHTTP(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.disabled {
			h.ServeHTTP(w, r)
			return
		}

		for _, rule := range s.rules {
			if !rule.Match(r) {
				// If no matches, just continue
				h.ServeHTTP(w, r)
				return
			}
		}

		s.plugins.Run(w, r, h)
	}
}
