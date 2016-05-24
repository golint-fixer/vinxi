package manager

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

type JSONLinks map[string]string

type JSONRule struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Config      config.Config `json:"config,omitempty"`
	Metadata    config.Config `json:"metadata,omitempty"`
	Links       JSONLinks     `json:"links,omitempty"`
}

type JSONPlugin struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Enabled     bool          `json:"enabled,omitempty"`
	Config      config.Config `json:"config,omitempty"`
	Metadata    config.Config `json:"metadata,omitempty"`
	Links       JSONLinks     `json:"links,omitempty"`
}

type JSONScope struct {
	ID      string       `json:"id"`
	Name    string       `json:"name,omitempty"`
	Rules   []JSONRule   `json:"rules,omitempty"`
	Plugins []JSONPlugin `json:"plugins,omitempty"`
	Links   JSONLinks    `json:"links,omitempty"`
}

type JSONInstance struct {
	ID          string          `json:"id"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Metadata    []config.Config `json:"metadata,omitempty"`
	Scopes      []JSONScope     `json:"scopes"`
	Links       JSONLinks       `json:"links,omitempty"`
}

type ControllerHandler func(http.ResponseWriter, *http.Request, *Controller)

type Controller struct {
	Path    string
	Method  string
	Manager *Manager
	Handler ControllerHandler
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler(w, r, c)
}

var routes = []*Controller{}

func AddRoute(method, path string, fn ControllerHandler) {
	route := &Controller{
		Path:    path,
		Method:  method,
		Handler: fn,
	}
	routes = append(routes, route)
}

func replyWithError(w http.ResponseWriter, status int, message string) {
	buf, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}{status, message})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(status)
	w.Write(buf)
}

func reply(w http.ResponseWriter) func([]byte, error) {
	return func(buf []byte, err error) {
		if err != nil {
			replyWithError(w, 500, err.Error())
		} else {
			w.Write(buf)
		}
	}
}

func notFound(w http.ResponseWriter, message string) {
	replyWithError(w, 404, message)
}

func init() {
	AddRoute("GET", "/", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		io.WriteString(w, "vinxi HTTP API manager "+vinxi.Version)
	})

	AddRoute("GET", "/version", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		reply(w)(json.Marshal(struct {
			Vinxi string `json:"vinxi"`
		}{
			Vinxi: vinxi.Version,
		}))
	})

	AddRoute("GET", "/health", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{}")
	})

	AddRoute("GET", "/catalog", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		links := struct {
			Links map[string]string `json:"links,omitempty"`
		}{
			Links: map[string]string{
				"self":    "/catalog",
				"plugins": "/catalog/plugins",
				"rules":   "/catalog/rules",
			},
		}

		reply(w)(json.Marshal(links))
	})

	AddRoute("GET", "/catalog/plugins", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		reply(w)(json.Marshal(plugin.Plugins))
	})

	AddRoute("GET", "/catalog/rules", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		reply(w)(json.Marshal(rule.Rules))
	})

	AddRoute("GET", "/plugins", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		plugins := createPlugins(c.Manager.Plugins.All())
		buf, err := json.Marshal(plugins)
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}
		w.Write(buf)
	})

	AddRoute("GET", "/scopes", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		scopes := createScopes(c.Manager.Scopes())
		buf, err := json.Marshal(scopes)
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}
		w.Write(buf)
	})

	AddRoute("GET", "/scopes/:scope", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":scope")

		scope := c.Manager.GetScope(id)
		if scope == nil {
			notFound(w, "Scope not found")
			return
		}

		data, err := encodeJSON(createScope(scope))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}
		w.Write(data)
	})

	AddRoute("GET", "/instances", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		buf, err := json.Marshal(c.Manager.instances)
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}
		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		links := JSONLinks{
			"self":   "/instances/" + id,
			"scopes": "/instances/" + id + "/scopes",
		}

		node := JSONInstance{
			ID:          instance.ID,
			Name:        instance.Name,
			Description: instance.Description,
			Links:       links,
		}

		node.Scopes = createScopes(instance.GetScopes())

		buf, err := json.Marshal(node)
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance/scopes", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		buf, err := json.Marshal(createScopes(instance.GetScopes()))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance/scopes/:scope", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")
		scopeId := req.URL.Query().Get(":scope")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		scope := instance.GetScope(scopeId)
		if scope == nil {
			notFound(w, "Scope not found")
			return
		}

		buf, err := json.Marshal(createScope(scope))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance/scopes/:scope/plugins", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")
		scopeId := req.URL.Query().Get(":scope")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		scope := instance.GetScope(scopeId)
		if scope == nil {
			notFound(w, "Scope not found")
			return
		}

		buf, err := json.Marshal(createPlugins(scope.Plugins.All()))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance/scopes/:scope/plugins/:plugin", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")
		scopeId := req.URL.Query().Get(":scope")
		pluginId := req.URL.Query().Get(":plugin")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		scope := instance.GetScope(scopeId)
		if scope == nil {
			notFound(w, "Scope not found")
			return
		}

		plugin := scope.Plugins.Get(pluginId)
		if plugin == nil {
			notFound(w, "Plugin not found")
			return
		}

		buf, err := json.Marshal(createPlugin(plugin))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance/scopes/:scope/rules", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")
		scopeId := req.URL.Query().Get(":scope")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		scope := instance.GetScope(scopeId)
		if scope == nil {
			notFound(w, "Scope not found")
			return
		}

		buf, err := json.Marshal(createRules(scope))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})

	AddRoute("GET", "/instances/:instance/scopes/:scope/rules/:rule", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":instance")
		scopeId := req.URL.Query().Get(":scope")
		ruleId := req.URL.Query().Get(":rule")

		instance := c.Manager.GetInstance(id)
		if instance == nil {
			notFound(w, "Instance not found")
			return
		}

		scope := instance.GetScope(scopeId)
		if scope == nil {
			notFound(w, "Scope not found")
			return
		}

		rule := scope.Rules.Get(ruleId)
		if rule == nil {
			notFound(w, "Plugin not found")
			return
		}

		buf, err := json.Marshal(createRule(rule))
		if err != nil {
			replyWithError(w, 500, err.Error())
			return
		}

		w.Write(buf)
	})
}

func encodeJSON(data interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	return buf.Bytes(), err
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
	buf := make([]JSONScope, len(scopes))
	for i, scope := range scopes {
		buf[i] = createScope(scope)
	}
	return buf
}

func createRules(scope *Scope) []JSONRule {
	rules := make([]JSONRule, scope.Rules.Len())
	for i, rule := range scope.Rules.All() {
		rules[i] = createRule(rule)
	}
	return rules
}

func createPlugins(plugins []plugin.Plugin) []JSONPlugin {
	list := make([]JSONPlugin, len(plugins))
	for i, plugin := range plugins {
		list[i] = createPlugin(plugin)
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
