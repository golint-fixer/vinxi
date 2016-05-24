package manager

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/h2non/httprouter"
	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/layer"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"

	// An empty import is required to load all the rules subpackages
	_ "gopkg.in/vinxi/vinxi.v0/plugins"
	_ "gopkg.in/vinxi/vinxi.v0/rules"
)

// Manager represents the vinxi proxy admin manager.
type Manager struct {
	layer     *layer.Layer
	plugins   *plugin.Layer
	sm        sync.RWMutex
	scopes    []*Scope
	im        sync.RWMutex
	instances []*Instance
	Server    *http.Server
	Router    *httprouter.Router
}

// New creates a new manager able to manage
// and configure multiple vinxi proxy instance.
func New() *Manager {
	return &Manager{layer: layer.New(), plugins: plugin.NewLayer(), Router: httprouter.New()}
}

// Manage creates a new empty manage and
// starts managing the given vinxi proxy instance.
func Manage(name, description string, proxy *vinxi.Vinxi) *Manager {
	m := New()
	m.Manage(name, description, proxy)
	return m
}

// Manage adds a new vinxi proxy instance to be
// managed by the current manager instance.
func (m *Manager) Manage(name, description string, proxy *vinxi.Vinxi) {
	// Register manager middleware in the proxy
	proxy.Layer.UsePriority(layer.RequestPhase, layer.Tail, m)

	// Register the managed Vinxi instance
	instance := NewInstance(name, description, proxy)

	m.im.Lock()
	m.instances = append(m.instances, instance)
	m.im.Unlock()
}

// Use attaches a new middleware handler for incoming HTTP traffic.
func (m *Manager) Use(handler ...interface{}) {
	m.layer.Use(layer.RequestPhase, handler...)
}

// UsePhase attaches a new middleware handler to a specific phase.
func (m *Manager) UsePhase(phase string, handler ...interface{}) {
	m.layer.Use(phase, handler...)
}

// UseFinalHandler uses a new middleware handler function as final handler.
func (m *Manager) UseFinalHandler(fn http.Handler) {
	m.layer.UseFinalHandler(fn)
}

// ListenAndServe creates a new admin HTTP server and starts listening on
// the network based on the given server options.
func (m *Manager) ListenAndServe(opts ServerOptions) (*http.Server, error) {
	m.Server = NewServer(opts)
	m.configure()
	return m.Server, Listen(m.Server, opts)
}

// ServeDefault creates a new admin HTTP server and starts listening
// on the network based on the default server settings.
func (m *Manager) ServeDefault() (*http.Server, error) {
	return m.ListenAndServe(ServerOptions{})
}

// NewScope creates a new scope based on the given name
// and optional description.
func (m *Manager) NewScope(name, description string) *Scope {
	scope := NewScope(name, description)
	m.sm.Lock()
	m.scopes = append(m.scopes, scope)
	m.sm.Unlock()
	return scope
}

// NewScope creates a new default scope.
func (m *Manager) NewDefaultScope(rules ...rule.Rule) *Scope {
	scope := m.NewScope("default", "Default generic scope")
	scope.UseRule(rules...)
	return scope
}

// GetScope finds and returns a vinxi managed instance.
func (m *Manager) GetScope(name string) *Scope {
	m.sm.Lock()
	defer m.sm.Unlock()

	for _, scope := range m.scopes {
		if scope.ID == name || scope.Name == name {
			return scope
		}
	}

	return nil
}

// RemoveScope removes a registered scope.
// Returns false if the scope cannot be found.
func (m *Manager) RemoveScope(name string) bool {
	m.sm.Lock()
	defer m.sm.Unlock()

	for i, scope := range m.scopes {
		if scope.ID == name || scope.Name == name {
			m.scopes = append(m.scopes[:i], m.scopes[i+1:]...)
			return true
		}
	}

	return false
}

// GetInstance finds and returns a vinxi managed instance.
func (m *Manager) GetInstance(name string) *Instance {
	m.im.Lock()
	defer m.im.Unlock()

	for _, instance := range m.instances {
		if instance.ID == name || instance.Name == name {
			return instance
		}
	}

	return nil
}

// RemoveInstance removes a registered vinxi instance.
// Returns false if the instance cannot be found.
func (m *Manager) RemoveInstance(name string) bool {
	m.im.Lock()
	defer m.im.Unlock()

	for i, instance := range m.instances {
		if instance.ID == name || instance.Name == name {
			m.instances = append(m.instances[:i], m.instances[i+1:]...)
			return true
		}
	}

	return false
}

// HandleHTTP is triggered by the vinxi middleware layer on incoming HTTP request.
func (m *Manager) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	next := h

	m.sm.RLock()
	for _, scope := range m.scopes {
		next = http.HandlerFunc(scope.HandleHTTP(next))
	}
	m.sm.RUnlock()

	next.ServeHTTP(w, r)
}

// serveHTTP is used to handle HTTP traffic via admin API.
func (m *Manager) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// trigger middleware layer first, then the router
	m.layer.Run(layer.RequestPhase, w, r, m.Router)
}

// configure is used to configure the HTTP API.
func (m *Manager) configure() {
	// bind the admin http.Handler
	m.Server.Handler = http.HandlerFunc(m.serveHTTP)

	// Define route handlers
	for _, r := range routes {
		r.Manager = m // Expose manager instance in routes
		m.Router.Handler(r.Method, r.Path, r)
	}
}

type JSONRule struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Config      config.Config `json:"config,omitempty"`
}

type JSONPlugin struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Enabled     bool          `json:"enabled,omitempty"`
	Config      config.Config `json:"config,omitempty"`
}

type JSONScope struct {
	ID      string       `json:"id"`
	Name    string       `json:"name,omitempty"`
	Rules   []JSONRule   `json:"rules,omitempty"`
	Plugins []JSONPlugin `json:"plugins,omitempty"`
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

func reply(w http.ResponseWriter) func([]byte, error) {
	return func(buf []byte, err error) {
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(buf)
	}
}

func init() {
	AddRoute("GET", "/", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		io.WriteString(w, "vinxi HTTP API manager "+vinxi.Version)
	})

	AddRoute("GET", "/version", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		versions := struct {
			Vinxi string `json:"vinxi"`
		}{
			Vinxi: vinxi.Version,
		}

		reply(w)(json.Marshal(versions))
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

	AddRoute("GET", "/instances", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		buf := &bytes.Buffer{}

		err := json.NewEncoder(buf).Encode(c.Manager.instances)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(buf.Bytes())
	})

	AddRoute("GET", "/instances/:instance", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		buf := &bytes.Buffer{}
		id := req.URL.Query().Get(":instance")
		// panic("adasdasdsad: " + id)

		mgr := c.Manager
		instance := mgr.GetInstance(id)
		if instance == nil {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
			return
		}

		scopes := createScopes(instance.GetScopes())
		err := json.NewEncoder(buf).Encode(scopes)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(buf.Bytes())
	})

	AddRoute("GET", "/plugins", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		buf := &bytes.Buffer{}
		scopes := createScopes(c.Manager.scopes)

		err := json.NewEncoder(buf).Encode(scopes)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(buf.Bytes())
	})

	AddRoute("GET", "/scopes", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		buf := &bytes.Buffer{}
		scopes := createScopes(c.Manager.scopes)

		err := json.NewEncoder(buf).Encode(scopes)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(buf.Bytes())
	})

	AddRoute("GET", "/scopes/:scope", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		id := req.URL.Query().Get(":scope")

		// Find scope by ID
		for _, scope := range c.Manager.scopes {
			if scope.ID == id {
				data, err := encodeJSON(createScope(scope))
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}
				w.Write(data)
				return
			}
		}

		w.WriteHeader(404)
		w.Write([]byte("not found"))
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
		Plugins: createPlugins(scope),
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
	for i, rule := range scope.Rules.Get() {
		rules[i] = JSONRule{ID: rule.ID(), Name: rule.Name(), Description: rule.Description(), Config: rule.Config()}
	}
	return rules
}

func createPlugins(scope *Scope) []JSONPlugin {
	plugins := make([]JSONPlugin, scope.Plugins.Len())
	for i, plugin := range scope.Plugins.Get() {
		plugins[i] = JSONPlugin{ID: plugin.ID(), Name: plugin.Name(), Description: plugin.Description(), Config: plugin.Config()}
	}
	return plugins
}
