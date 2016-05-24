package manager

import (
	"net/http"
	"sync"

	"github.com/dchest/uniuri"
	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// Instance represents the manager instance level.
type Instance struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	sm          sync.RWMutex
	scopes      []*Scope
	instance    *vinxi.Vinxi
}

// NewInstance creates a new vinxi manager instance.
func NewInstance(name, description string, proxy *vinxi.Vinxi) *Instance {
	return &Instance{ID: uniuri.New(), Name: name, Description: description, instance: proxy}
}

// NewScope creates a new scope based on the given name
// and optional description.
func (i *Instance) NewScope(name, description string) *Scope {
	scope := NewScope(name, description)
	i.sm.Lock()
	i.scopes = append(i.scopes, scope)
	i.sm.Unlock()
	return scope
}

// NewScope creates a new default scope.
func (i *Instance) NewDefaultScope(rules ...rule.Rule) *Scope {
	scope := i.NewScope("default", "Default generic scope")
	scope.UseRule(rules...)
	return scope
}

// Scopes returns the list of registered scopes.
func (i *Instance) Scopes() []*Scope {
	i.sm.RLock()
	defer i.sm.RUnlock()
	return i.scopes
}

// GetScope finds and return a registered scope instance.
func (i *Instance) GetScope(name string) *Scope {
	i.sm.Lock()
	defer i.sm.Unlock()

	for _, scope := range i.scopes {
		if scope.ID == name || scope.Name == name {
			return scope
		}
	}

	return nil
}

// RemoveScope removes a registered scope.
// Returns false if the scope cannot be found.
func (i *Instance) RemoveScope(name string) bool {
	i.sm.Lock()
	defer i.sm.Unlock()

	for x, scope := range i.scopes {
		if scope.ID == name || scope.Name == name {
			i.scopes = append(i.scopes[:x], i.scopes[x+1:]...)
			return true
		}
	}

	return false
}

// SetMeta defines a new metadata metadata field.
// TODO protect by mutex
func (i *Instance) SetMeta(key string, value interface{}) {
	i.Metadata[key] = value
}

// HandleHTTP is triggered by the vinxi middleware layer on incoming HTTP request.
func (i *Instance) HandleHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	i.sm.RLock()
	for _, scope := range i.scopes {
		next = http.HandlerFunc(scope.HandleHTTP(next))
	}
	i.sm.RUnlock()

	next.ServeHTTP(w, r)
}
