package plugin

import (
	"net/http"
	"sync"
)

// Layer represents a plugins layer designed to intrument
// proxies providing plugin based dynamic configuration
// capabilities, such as register/unregister or
// enable/disable plugins at runtime satefy.
type Layer struct {
	rwm  sync.RWMutex
	pool []Plugin
}

// NewLayer creates a new plugins layer.
func NewLayer() *Layer {
	return &Layer{}
}

// Use registers one or multiples plugins in the current plugin layer.
func (l *Layer) Use(plugin ...Plugin) {
	l.rwm.Lock()
	l.pool = append(l.pool, plugin...)
	l.rwm.Unlock()
}

// Len returns the registered plugins length.
func (l *Layer) Len() int {
	return len(l.pool)
}

// Get returns an slice of the registered plugins.
func (l *Layer) Get() []Plugin {
	l.rwm.Lock()
	defer l.rwm.Unlock()
	return l.pool
}

// Remove removes a plugin looking by its unique identifier.
func (l *Layer) Remove(id string) bool {
	l.rwm.Lock()
	defer l.rwm.Unlock()

	for i, plugin := range l.pool {
		if plugin.ID() == id {
			l.pool = append(l.pool[:i], l.pool[i+1:]...)
			return true
		}
	}

	return false
}

// Run triggers the plugins layer call chain.
// This function is designed to be executed by top-level middleware layers.
func (l *Layer) Run(w http.ResponseWriter, r *http.Request, h http.Handler) {
	next := h
	l.rwm.RLock()
	for _, plugin := range l.pool {
		next = plugin.HandleHTTP(next)
	}
	l.rwm.RUnlock()
	next.ServeHTTP(w, r)
}
