package rule

import (
	"net/http"
	"sync"
)

// Layer represents a rules layer designed to intrument
// proxies providing plugin based dynamic configuration
// capabilities, such as register/unregister or
// enable/disable rules at runtime satefy.
type Layer struct {
	rwm  sync.RWMutex
	pool []Rule
}

// NewLayer creates a new rules layer.
func NewLayer() *Layer {
	return &Layer{}
}

// Use registers one or multiples plugins in the current rule layer.
func (l *Layer) Use(rule ...Rule) {
	l.rwm.Lock()
	l.pool = append(l.pool, rule...)
	l.rwm.Unlock()
}

// Len returns the registered rules length.
func (l *Layer) Len() int {
	l.rwm.RLock()
	defer l.rwm.RUnlock()
	return len(l.pool)
}

// Remove removes a rule looking by its unique identifier.
func (l *Layer) Remove(id string) bool {
	l.rwm.Lock()
	defer l.rwm.Unlock()

	for i, rule := range l.pool {
		if rule.ID() == id {
			l.pool = append(l.pool[:i], l.pool[i+1:]...)
			return true
		}
	}

	return false
}

// Get finds and treturns a rule instance.
func (l *Layer) Get(name string) Rule {
	l.rwm.Lock()
	defer l.rwm.Unlock()

	for _, rule := range l.pool {
		if rule.ID() == name || rule.Name() == name {
			return rule
		}
	}
	return nil
}

// All returns the slice of registered rules.
func (l *Layer) All() []Rule {
	l.rwm.Lock()
	defer l.rwm.Unlock()
	return l.pool
}

// Flush removes all the registered rules.
func (l *Layer) Flush() {
	l.rwm.Lock()
	l.pool = []Rule{}
	l.rwm.Unlock()
}

// Match matches the given http.Request agains the registered rules.
// If all the rules passes it will return true, otherwise false.
func (l *Layer) Match(r *http.Request) bool {
	l.rwm.RLock()
	defer l.rwm.RUnlock()

	for _, rule := range l.pool {
		if !rule.Match(r) {
			return false
		}
	}

	return true
}
