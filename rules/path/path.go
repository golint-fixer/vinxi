package path

import (
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/mux"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

const (
	// Name exposes the rule name identifier.
	Name = "path"
	// Description exposes the rule semantic description.
	Description = "Matches HTTP request URL path againts a given path pattern"
)

// params defines the rule specific configuration params.
var params = rule.Params{
	rule.Field{
		Name:        "path",
		Type:        "string",
		Description: "Path expression to match",
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
	return rule.Matcher(mux.MatchPath(opts.GetString("path"))), nil
}

// New creates a new rule who filters the traffic
// if matches with the following path expression.
// Regular expressions is supported.
func New(path string) (rule.Rule, error) {
	return rule.NewWithConfig(Rule, config.Config{"path": path})
}

func init() {
	rule.Register(Rule)
}
