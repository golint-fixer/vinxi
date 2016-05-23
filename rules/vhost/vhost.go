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
	Factory:     Factory,
	Params:      params,
}

// Factory represents the rule factory function
// designed to be called via rules constructor.
func Factory(opts config.Config) rule.Rule {
	return rule.NewRuleWithConfig(
		Name,
		Description,
		opts,
		mux.MatchHost(opts.GetString("host")),
	)
}

// New creates a new rule who filters the traffic
// if matches with the following path expression.
// Regular expressions is supported.
func New(host string) rule.Rule {
	config := map[string]interface{}{"host": host}
	return Factory(config)
}

func init() {
	rule.Register(Rule)
}
