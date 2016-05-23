package manager

import (
	"net/http"

	"github.com/dchest/uniuri"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// Scope represents the HTTP configuration scope who can
// store rules and plugins.
type Scope struct {
	// rules stores the scope registered rules.
	rules *rule.Layer
	// plugins provides the plugin register layer.
	plugins *plugin.Layer
	// ID is used to store the plugin unique identifier.
	ID string
	// Name is used to store the scope semantic alias.
	Name string
	// Description is used to store the scope human
	// friendly description.
	Description string
}

// NewScope creates a new Scope instance
// with the given name alias and optional description.
func NewScope(name, description string) *Scope {
	return &Scope{
		ID:          uniuri.New(),
		Name:        name,
		Description: description,
		rules:       rule.NewLayer(),
		plugins:     plugin.NewLayer(),
	}
}

// UseRule registers one or multiple rules in the current scope.
func (s *Scope) UseRule(rules ...rule.Rule) {
	s.rules.Use(rules...)
}

// UseRule registers one or multiple plugins in the current scope.
func (s *Scope) UsePlugin(plugins ...plugin.Plugin) {
	s.plugins.Use(plugins...)
}

// Rules returns the rules register layer of the current scope.
func (s *Scope) Rules() *rule.Layer {
	return s.rules
}

// Plugins returns the plugins register layer of the current scope.
func (s *Scope) Plugins() *plugin.Layer {
	return s.plugins
}

// RemoveRule removes a rule by its ID.
func (s *Scope) RemoveRule(id string) bool {
	return s.plugins.Remove(id)
}

// RemovePlugin removes a plugin by its ID.
func (s *Scope) RemovePlugin(id string) bool {
	return s.plugins.Remove(id)
}

// HandleHTTP is used to trigger the scope layer.
// If all the rules passes, it will execute the scope specific registered plugins.
func (s *Scope) HandleHTTP(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for !s.rules.Match(r) {
			// If no matches, just continue
			h.ServeHTTP(w, r)
			return
		}
		s.plugins.Run(w, r, h)
	}
}
