package manager

import (
	"net/http"
	"sync"

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
	// Config is used to retrieve the rule config.
	Config() Config
	// Match is used to determine if a given http.Request
	// passes the rule assertion.
	Match(*http.Request) bool
}

type matcher func(*http.Request) bool

type rule struct {
	id, name, description string
	matcher               matcher
	config                Config
}

func (r *rule) ID() string {
	return r.id
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) Description() string {
	return r.description
}

func (r *rule) Config() Config {
	return r.config
}

func (r *rule) Match(req *http.Request) bool {
	return r.matcher(req)
}

// NewRule creates a new rule entity based on the given matcher function.
func NewRule(name, description string, matcher func(*http.Request) bool) *rule {
	return NewRuleWithConfig(name, description, make(map[string]interface{}), matcher)
}

// NewRuleWithConfig creates a new rule entity based on the given config and matcher function.
func NewRuleWithConfig(name, description string, config map[string]interface{}, matcher func(*http.Request) bool) *rule {
	return &rule{
		id:          uniuri.New(),
		name:        name,
		description: description,
		matcher:     matcher,
		config:      Config(config),
	}
}

// RuleLayer represents a rules layer designed to intrument
// proxies providing plugin based dynamic configuration
// capabilities, such as register/unregister or
// enable/disable rules at runtime satefy.
type RuleLayer struct {
	rwm  sync.RWMutex
	pool []Rule
}

// NewRuleLayer creates a new rules layer.
func NewRuleLayer() *RuleLayer {
	return &RuleLayer{}
}

// Use registers one or multiples plugins in the current rule layer.
func (l *RuleLayer) Use(rule ...Rule) {
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
