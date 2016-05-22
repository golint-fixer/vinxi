package rules

import (
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/mux"
)

const (
	pathRuleName        = "path"
	pathRuleDescription = "Matches HTTP request URL path againts a given path pattern"
)

var pathConfig = Config{
	ConfigType{
		Name:        "path",
		Type:        "string",
		Description: "Path expression to match",
		Mandatory:   true,
	},
}

var pathRule = Rule{
	Name:        pathRuleName,
	Description: pathRuleDescription,
	Factory:     pathFactory,
	Config:      pathConfig,
}

func pathFactory(config config.Config) Rule {
	return NewRuleWithConfig(
		pathRuleName,
		pathRuleDescription,
		config,
		mux.MatchPath(config.GetString("path")),
	)
}

// Path creates a new rule who filters the traffic
// if matches with the following path expression.
// Regular expressions is supported.
func Path(path string) Rule {
	config := map[string]interface{}{"path": path}
	return pathFactory(config)
}

func init() {
	Rules.Register(pathRule)
}
