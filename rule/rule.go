package rule

import (
	"net/http"

	"github.com/dchest/uniuri"
	"gopkg.in/vinxi/vinxi.v0/config"
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
	Config() config.Config
	// Match is used to determine if a given http.Request
	// passes the rule assertion.
	Match(*http.Request) bool
}

type matcher func(*http.Request) bool

type rule struct {
	id, name, description string
	matcher               matcher
	config                config.Config
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

func (r *rule) Config() config.Config {
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
func NewRuleWithConfig(name, description string, opts map[string]interface{}, matcher func(*http.Request) bool) *rule {
	return &rule{
		id:          uniuri.New(),
		name:        name,
		description: description,
		matcher:     matcher,
		config:      config.Config(opts),
	}
}
