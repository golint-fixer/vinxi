package rule

import (
	"net/http"

	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/utils"
)

// Matcher represents the matching function interface
// used by rules to determine if a given traffic
// should be filtered or not.
type Matcher func(*http.Request) bool

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

// rule implements the Rule interface.
type rule struct {
	id, name, description string
	matcher               Matcher
	config                config.Config
}

// ID is used to retrieve the rule unique identifier.
func (r *rule) ID() string {
	return r.id
}

// Name is used to retrieve the rule semantic identifier.
func (r *rule) Name() string {
	return r.name
}

// Description is used to retrieve the rule human friendly description.
func (r *rule) Description() string {
	return r.description
}

// Config is used to retrieve the rule configuration params.
func (r *rule) Config() config.Config {
	return r.config
}

// Match is executed by the scope layer when a new request is recieved by the proxy.
func (r *rule) Match(req *http.Request) bool {
	return r.matcher(req)
}

// New creates a new rule entity based on the given matcher function.
func New(info Info) func(config.Config) (Rule, error) {
	return func(opts config.Config) (Rule, error) {
		return NewWithConfig(info, opts)
	}
}

// NewWithConfig creates a new rule entity based on the given config and matcher function.
func NewWithConfig(info Info, opts config.Config) (Rule, error) {
	if err := Validate(info.Params, opts); err != nil {
		return nil, err
	}

	// Build the rule matcher function
	matcher, err := info.Factory(opts)
	if err != nil {
		return nil, err
	}

	return &rule{
		id:          utils.NewID(),
		name:        info.Name,
		description: info.Description,
		config:      opts,
		matcher:     matcher,
	}, nil
}
