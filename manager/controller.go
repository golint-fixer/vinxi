package manager

import (
	"encoding/json"
	"net/http"

	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// ControllerHandler represents HTTP controller handler function
// interface used in routes.
type ControllerHandler func(*Context)

// Context is used to share request context entities
// across controllers.
type Context struct {
	Manager  *Manager
	Scope    *Scope
	Instance *Instance
	Request  *http.Request
	Response http.ResponseWriter
	Rule     rule.Rule
	Plugin   plugin.Plugin
}

// SendJSON is used to serialize and write the response as JSON.
func (c *Context) SendJSON(data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		c.SendError(500, err.Error())
		return
	}
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.Write(buf)
}

// Error replies with an custom error message and 500 as status code.
func (c *Context) SendError(status int, message string) {
	c.Response.Header().Set("Content-Type", "application/json")

	buf, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}{status, message})

	if err != nil {
		c.Response.WriteHeader(500)
		c.Response.Write([]byte(err.Error()))
		return
	}

	c.Response.WriteHeader(status)
	c.Response.Write(buf)
}

// SendNoContent replies with 204 status code.
func (c *Context) SendNoContent() {
	c.Response.WriteHeader(204)
}

// SendNoContent replies with 204 status code.
func (c *Context) SendNotFound(message string) {
	c.SendError(404, message)
}

// Controller represents a route handler.
type Controller struct {
	Path    string
	Method  string
	Manager *Manager
	Handler ControllerHandler
}

// ServeHTTP implements the http.HandlerFunc interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{Manager: c.Manager, Request: r, Response: w}
	c.handle(ctx)
}

func (c *Controller) handle(ctx *Context) {
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
		if ctx.Scope != nil {
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
