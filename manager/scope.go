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
	// Rules stores the scope registered rules.
	Rules *rule.Layer
	// Plugins provides the plugin register layer.
	Plugins *plugin.Layer
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
		Rules:       rule.NewLayer(),
		Plugins:     plugin.NewLayer(),
	}
}

// UseRule registers one or multiple rules in the current scope.
func (s *Scope) UseRule(rules ...rule.Rule) {
	s.Rules.Use(rules...)
}

// UsePlugin registers one or multiple plugins in the current scope.
func (s *Scope) UsePlugin(plugins ...plugin.Plugin) {
	s.Plugins.Use(plugins...)
}

// RemoveRule removes a rule by its ID.
func (s *Scope) RemoveRule(id string) bool {
	return s.Rules.Remove(id)
}

// FlushRules removes all the registered rules.
func (s *Scope) FlushRules() {
	s.Rules.Flush()
}

// RemovePlugin removes a plugin by its ID.
func (s *Scope) RemovePlugin(id string) bool {
	return s.Plugins.Remove(id)
}

// FlushPlugins removes all the registered plugins.
func (s *Scope) FlushPlugins() {
	s.Plugins.Flush()
}

// HandleHTTP is used to trigger the scope layer.
// If all the rules passes, it will execute the scope specific registered plugins.
func (s *Scope) HandleHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for !s.Rules.Match(r) {
			// If no matches, just continue
			h.ServeHTTP(w, r)
			return
		}
		s.Plugins.HandleHTTP(w, r, h)
	})
}
