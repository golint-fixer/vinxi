package manager

import (
	"os"
	"runtime"

	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// routes stores the registered routes.
var routes = []*Controller{}

func addRoute(method, path string, fn ControllerHandler) {
	route := &Controller{
		Path:    path,
		Method:  method,
		Handler: fn,
	}
	routes = append(routes, route)
}

// JSONRule represents the Rule entity for JSON serialization.
type JSONRule struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Config      config.Config `json:"config,omitempty"`
	Metadata    config.Config `json:"metadata,omitempty"`
}

// JSONPlugin represents the Plugin entity for JSON serialization.
type JSONPlugin struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Enabled     bool          `json:"enabled,omitempty"`
	Config      config.Config `json:"config,omitempty"`
	Metadata    config.Config `json:"metadata,omitempty"`
}

// JSONScope represents the Scope entity for JSON serialization.
type JSONScope struct {
	ID      string       `json:"id"`
	Name    string       `json:"name,omitempty"`
	Rules   []JSONRule   `json:"rules"`
	Plugins []JSONPlugin `json:"plugins"`
}

// JSONInstance represents the Instance entity for JSON serialization.
type JSONInstance struct {
	ID          string          `json:"id"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Metadata    []config.Config `json:"metadata,omitempty"`
	Scopes      []JSONScope     `json:"scopes"`
}

func init() {
	addRoute("GET", "/", func(ctx *Context) {
		hostname, _ := os.Hostname()

		info := struct {
			Hostname string            `json:"hostname"`
			Version  string            `json:"version"`
			Platform string            `json:"platform"`
			Links    map[string]string `json:"links"`
		}{
			Hostname: hostname,
			Version:  vinxi.Version,
			Platform: runtime.GOOS,
			Links: map[string]string{
				"catalog":   "/catalog",
				"plugins":   "/plugins",
				"scopes":    "/scopes",
				"instances": "/instances",
			},
		}

		ctx.SendJSON(info)
	})

	addRoute("GET", "/catalog", func(ctx *Context) {
		rules := []rule.Info{}
		for _, rule := range rule.Rules {
			rules = append(rules, rule)
		}

		plugins := []plugin.Info{}
		for _, plugin := range plugin.Plugins {
			plugins = append(plugins, plugin)
		}

		catalog := struct {
			Rules   []rule.Info   `json:"rules"`
			Plugins []plugin.Info `json:"plugins"`
		}{
			Rules:   rules,
			Plugins: plugins,
		}

		ctx.SendJSON(catalog)
	})

	addRoute("GET", "/plugins", func(ctx *Context) {
		ctx.SendJSON(createPlugins(ctx.Manager.Plugins.All()))
	})

	addRoute("GET", "/plugins/:plugin", func(ctx *Context) {
		ctx.SendJSON(createPlugin(ctx.Plugin))
	})

	addRoute("GET", "/scopes", func(ctx *Context) {
		ctx.SendJSON(createScopes(ctx.Manager.Scopes()))
	})

	addRoute("GET", "/scopes/:scope", func(ctx *Context) {
		ctx.SendJSON(createScope(ctx.Scope))
	})

	addRoute("GET", "/instances", func(ctx *Context) {
		ctx.SendJSON(createInstances(ctx.Manager.Instances()))
	})

	addRoute("GET", "/instances/:instance", func(ctx *Context) {
		ctx.SendJSON(createInstance(ctx.Instance))
	})

	addRoute("DELETE", "/instances/:instance", func(ctx *Context) {
		if ctx.Manager.RemoveInstance(ctx.Instance.ID) {
			ctx.SendNoContent()
		} else {
			ctx.SendError(500, "Cannot remove instance")
		}
	})

	addRoute("GET", "/instances/:instance/scopes", func(ctx *Context) {
		ctx.SendJSON(createScopes(ctx.Instance.Scopes()))
	})

	addRoute("GET", "/instances/:instance/scopes/:scope", func(ctx *Context) {
		ctx.SendJSON(createScope(ctx.Scope))
	})

	addRoute("DELETE", "/instances/:instance/scopes/:scope", func(ctx *Context) {
		if ctx.Instance.RemoveScope(ctx.Scope.ID) {
			ctx.SendNoContent()
		} else {
			ctx.SendError(500, "Cannot remove scope")
		}
	})

	addRoute("GET", "/instances/:instance/scopes/:scope/plugins", func(ctx *Context) {
		ctx.SendJSON(createPlugins(ctx.Scope.Plugins.All()))
	})

	addRoute("GET", "/instances/:instance/scopes/:scope/plugins/:plugin", func(ctx *Context) {
		ctx.SendJSON(createPlugin(ctx.Plugin))
	})

	addRoute("DELETE", "/instances/:instance/scopes/:scope/plugins/:plugin", func(ctx *Context) {
		if ctx.Scope.RemovePlugin(ctx.Plugin.ID()) {
			ctx.SendNoContent()
		} else {
			ctx.SendError(500, "Cannot remove plugin")
		}
	})

	addRoute("GET", "/instances/:instance/scopes/:scope/rules", func(ctx *Context) {
		ctx.SendJSON(createRules(ctx.Scope))
	})

	addRoute("GET", "/instances/:instance/scopes/:scope/rules/:rule", func(ctx *Context) {
		ctx.SendJSON(createRule(ctx.Rule))
	})

	addRoute("DELETE", "/instances/:instance/scopes/:scope/rules/:rule", func(ctx *Context) {
		if ctx.Scope.RemoveRule(ctx.Rule.ID()) {
			ctx.SendNoContent()
		} else {
			ctx.SendError(500, "Cannot remove rule")
		}
	})
}

func createInstance(instance *Instance) JSONInstance {
	return JSONInstance{
		ID:          instance.ID,
		Name:        instance.Name,
		Description: instance.Description,
		Scopes:      createScopes(instance.Scopes()),
	}
}

func createInstances(instances []*Instance) []JSONInstance {
	list := []JSONInstance{}
	for _, instance := range instances {
		list = append(list, createInstance(instance))
	}
	return list
}

func createScope(scope *Scope) JSONScope {
	return JSONScope{
		ID:      scope.ID,
		Name:    scope.Name,
		Rules:   createRules(scope),
		Plugins: createPlugins(scope.Plugins.All()),
	}
}

func createScopes(scopes []*Scope) []JSONScope {
	list := []JSONScope{}
	for _, scope := range scopes {
		list = append(list, createScope(scope))
	}
	return list
}

func createRules(scope *Scope) []JSONRule {
	rules := []JSONRule{}
	for _, rule := range scope.Rules.All() {
		rules = append(rules, createRule(rule))
	}
	return rules
}

func createPlugins(plugins []plugin.Plugin) []JSONPlugin {
	list := []JSONPlugin{}
	for _, plugin := range plugins {
		list = append(list, createPlugin(plugin))
	}
	return list
}

func createRule(rule rule.Rule) JSONRule {
	return JSONRule{
		ID:          rule.ID(),
		Name:        rule.Name(),
		Description: rule.Description(),
		Config:      rule.Config(),
	}
}

func createPlugin(plugin plugin.Plugin) JSONPlugin {
	return JSONPlugin{
		ID:          plugin.ID(),
		Name:        plugin.Name(),
		Description: plugin.Description(),
		Config:      plugin.Config(),
		Metadata:    plugin.Metadata(),
	}
}
