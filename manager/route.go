package manager

import (
	"net/http"
	"regexp"
)

// RouteHandler represents HTTP router handler function
// interface used in routes.
type RouteHandler func(*Context)

// Route represents a route handler.
type Route struct {
	Path    string
	Method  string
	Manager *Manager
	Handler RouteHandler
}

// ServeHTTP implements the http.HandlerFunc interface.
func (c *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{Manager: c.Manager, Request: r, Response: w}
	c.handle(ctx)
}

func (c *Route) isManagerPath(ctx *Context) bool {
	matches, _ := regexp.MatchString("^/manager/", ctx.Request.URL.Path)
	return matches
}

func (c *Route) handle(ctx *Context) {
	if c.isManagerPath(ctx) {
		ctx.AdminPlugins = ctx.Manager.AdminPlugins
	}

	instanceId := ctx.Request.URL.Query().Get(":instance")
	if instanceId != "" {
		ctx.Instance = ctx.Manager.GetInstance(instanceId)
		if ctx.Instance == nil {
			ctx.SendNotFound("Instance not found")
			return
		}
	}

	scopeId := ctx.Request.URL.Query().Get(":scope")
	if scopeId != "" {
		if ctx.Instance != nil {
			ctx.Scope = ctx.Instance.GetScope(scopeId)
		} else {
			ctx.Scope = ctx.Manager.GetScope(scopeId)
		}
		if ctx.Scope == nil {
			ctx.SendNotFound("Scope not found")
			return
		}
	}

	pluginId := ctx.Request.URL.Query().Get(":plugin")
	if pluginId != "" {
		if ctx.AdminPlugins != nil {
			ctx.Plugin = ctx.AdminPlugins.Get(pluginId)
		} else if ctx.Scope != nil {
			ctx.Plugin = ctx.Scope.Plugins.Get(pluginId)
		} else {
			ctx.Plugin = ctx.Manager.GetPlugin(pluginId)
		}
		if ctx.Plugin == nil {
			ctx.SendNotFound("Plugin not found")
			return
		}
	}

	ruleId := ctx.Request.URL.Query().Get(":rule")
	if ruleId != "" && ctx.Scope != nil {
		ctx.Rule = ctx.Scope.Rules.Get(ruleId)
		if ctx.Rule == nil {
			ctx.SendNotFound("Rule not found")
			return
		}
	}

	// Finally run the router if all path validations are ok
	c.Handler(ctx)
}
