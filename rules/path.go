package rules

import (
	"gopkg.in/vinxi/vinxi.v0/manager"
	"gopkg.in/vinxi/vinxi.v0/mux"
)

const (
	pathRuleName        = "path"
	pathRuleDescription = "Matches HTTP request URL path againts a given path pattern"
)

var pathRule = Rule{
	Name:        pathRuleName,
	Description: pathRuleDescription,
	Factory:     pathFactory,
}

func pathFactory(config manager.Config) manager.Rule {
	return manager.NewRuleWithConfig(
		pathRuleName,
		pathRuleDescription,
		config,
		mux.MatchPath(config.GetString("path")),
	)
}

// Path creates a new rule who filters the traffic
// if matches with the following path expression.
// Regular expressions is supported.
func Path(path string) manager.Rule {
	config := map[string]interface{}{"path": path}
	return pathFactory(config)
}

func init() {
	Rules.Register(pathRule)
}
