package manager

import (
	"net/http"
	"sync"

	"gopkg.in/vinxi/vinxi.v0/rule"
)

// RuleLayer represents a rules layer designed to intrument
// proxies providing plugin based dynamic configuration
// capabilities, such as register/unregister or
// enable/disable rules at runtime satefy.
type RuleLayer struct {
	rwm  sync.RWMutex
	pool []rule.Rule
}

// NewRuleLayer creates a new rules layer.
func NewRuleLayer() *RuleLayer {
	return &RuleLayer{}
}

// Use registers one or multiples plugins in the current rule layer.
func (l *RuleLayer) Use(rule ...rule.Rule) {
	l.rwm.Lock()
	l.pool = append(l.pool, rule...)
	l.rwm.Unlock()
}

// Len returns the registered rules length.
func (l *RuleLayer) Len() int {
	l.rwm.RLock()
	defer l.rwm.RUnlock()
	return len(l.pool)
}

// Remove removes a rule looking by its unique identifier.
func (l *RuleLayer) Remove(id string) bool {
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

// Match matches the given http.Request agains the registered rules.
// If all the rules passes it will return true, otherwise false.
func (l *RuleLayer) Match(r *http.Request) bool {
	l.rwm.RLock()
	defer l.rwm.RUnlock()

	for _, rule := range l.pool {
		if !rule.Match(r) {
			return false
		}
	}

	return true
}
