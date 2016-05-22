package path

import (
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/mux"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

const (
	pathRuleName        = "path"
	pathRuleDescription = "Matches HTTP request URL path againts a given path pattern"
)

var pathParams = rule.Params{
	rule.Field{
		Name:        "path",
		Type:        "string",
		Description: "Path expression to match",
		Mandatory:   true,
	},
}

var pathRule = rule.Info{
	Name:        pathRuleName,
	Description: pathRuleDescription,
	Factory:     pathFactory,
	Params:      pathParams,
}

func pathFactory(opts config.Config) rule.Rule {
	return rule.NewRuleWithConfig(
		pathRuleName,
		pathRuleDescription,
		opts,
		mux.MatchPath(opts.GetString("path")),
	)
}

// Path creates a new rule who filters the traffic
// if matches with the following path expression.
// Regular expressions is supported.
func Path(path string) rule.Rule {
	config := map[string]interface{}{"path": path}
	return pathFactory(config)
}

func init() {
	rule.Register(pathRule)
}
