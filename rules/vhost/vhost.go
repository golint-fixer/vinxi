package vhost

import (
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/mux"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

const (
	// Name exposes the rule name identifier.
	Name = "vhost"
	// Description exposes the rule semantic description.
	Description = "Matches HTTP request by Host header"
)

// params defines the rule specific configuration params.
var params = rule.Params{
	rule.Field{
		Name:        "host",
		Type:        "string",
		Description: "Hostname expression to match. Regular expressions are supported.",
		Mandatory:   true,
	},
}

// Rule exposes the rule metadata information.
// Mostly used internally.
var Rule = rule.Info{
	Name:        Name,
	Description: Description,
	Factory:     factory,
	Params:      params,
}

// factory represents the rule factory function
// designed to be called via rules constructor.
func factory(opts config.Config) (rule.Matcher, error) {
	return rule.Matcher(mux.MatchHost(opts.GetString("host"))), nil
}

// New creates a new rule who filters the traffic
// if matches with the following path expression.
// Regular expressions is supported.
func New(host string) (rule.Rule, error) {
	return rule.NewWithConfig(Rule, config.Config{"host": host})
}

func init() {
	rule.Register(Rule)
}
