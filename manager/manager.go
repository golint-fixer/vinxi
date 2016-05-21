package manager

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/layer"
)

type Layer struct {
	layer *layer.Layer
}

type Rule interface {
	ID() string
	Name() string
	Description() string
	JSONConfig() string
	Match(*http.Request) bool
}

type Scope struct {
	disabled bool
	rules    []Rule
	plugins  *PluginLayer

	ID          string
	Name        string
	Description string
}

func NewScope(rules ...Rule) *Scope {
	return &Scope{ID: uniuri.New(), Name: "default", plugins: NewPluginLayer(), rules: rules}
}

func (s *Scope) UsePlugin(plugin Plugin) {
	s.plugins.Use(plugin)
}

func (s *Scope) AddRule(rules ...Rule) {
	s.rules = append(s.rules, rules...)
}

func (s *Scope) Rules() []Rule {
	return s.rules
}

func (s *Scope) Disable() {
	s.disabled = true
}

func (s *Scope) Enable() {
	s.disabled = false
}

func (s *Scope) HandleHTTP(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.disabled {
			h.ServeHTTP(w, r)
			return
		}

		for _, rule := range s.rules {
			if !rule.Match(r) {
				// Continue
				h.ServeHTTP(w, r)
				return
			}
		}

		s.plugins.Run(w, r, h)
	}
}

type VinxiInstance struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	scopes      []*Scope
	instance    *vinxi.Vinxi
}

type Manager struct {
	Server    *http.Server
	plugins   *PluginLayer
	scopes    []*Scope
	instances []*VinxiInstance
}

func New() *Manager {
	return &Manager{plugins: NewPluginLayer()}
}

func Manage(name, description string, proxy *vinxi.Vinxi) *Manager {
	m := New()
	m.Manage(name, description, proxy)
	return m
}

func (m *Manager) Manage(name, description string, proxy *vinxi.Vinxi) {
	// Register manager middleware in the proxy
	proxy.Layer.UsePriority("request", layer.Tail, m)

	// Register the managed Vinxi instance
	instance := &VinxiInstance{ID: uniuri.New(), Name: name, Description: description, instance: proxy}
	m.instances = append(m.instances, instance)
}

// ListenAndServe creates a new admin HTTP server and starts listening on
// the network based on the given server options.
func (m *Manager) ListenAndServe(opts ServerOptions) (*http.Server, error) {
	m.Server = NewServer(opts)
	m.Configure()
	return m.Server, Listen(m.Server, opts)
}

// ServeDefault creates a new admin HTTP server and starts listening
// on the network based on the default server settings.
func (m *Manager) ServeDefault() (*http.Server, error) {
	return m.ListenAndServe(ServerOptions{})
}

// NewScope creates a new
func (m *Manager) NewScope(rules ...Rule) *Scope {
	scope := NewScope(rules...)
	m.scopes = append(m.scopes, scope)
	return scope
}

func (m *Manager) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	next := h

	for _, scope := range m.scopes {
		next = http.HandlerFunc(scope.HandleHTTP(next))
	}

	next.ServeHTTP(w, r)
}

type JSONRule struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Config      string `json:"config,omitempty"`
}

type JSONPlugin struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
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

func init() {
	AddRoute("GET", "/", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		io.WriteString(w, "vinxi HTTP API manager "+vinxi.Version)
	})

	AddRoute("GET", "/catalog", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		io.WriteString(w, "Catalog here...")
	})

	AddRoute("GET", "/catalog/plugins", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		io.WriteString(w, "Plugin catalog here...")
	})

	AddRoute("GET", "/catalog/scopes", func(w http.ResponseWriter, req *http.Request, c *Controller) {
		io.WriteString(w, "Scopes catalog here...")
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

		mgr := c.Manager
		for _, instance := range mgr.instances {
			if instance.ID == id || instance.Name == id {
				scopes := createScopes(instance.scopes)

				err := json.NewEncoder(buf).Encode(scopes)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}

				w.Write(buf.Bytes())
				return
			}
		}

		w.WriteHeader(404)
		w.Write([]byte("Not found"))
		return
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

func (m *Manager) Configure() error {
	router := httprouter.New()
	m.Server.Handler = router

	// Define route handlers
	for _, r := range routes {
		r.Manager = m // Expose manager instance via routes
		router.Handler(r.Method, r.Path, r)
	}

	return nil
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
	rules := make([]JSONRule, len(scope.rules))
	for i, rule := range scope.rules {
		rules[i] = JSONRule{ID: rule.ID(), Name: rule.Name(), Description: rule.Description(), Config: rule.JSONConfig()}
	}
	return rules
}

func createPlugins(scope *Scope) []JSONPlugin {
	plugins := make([]JSONPlugin, scope.plugins.Len())
	for i, plugin := range scope.plugins.pool {
		plugins[i] = JSONPlugin{ID: plugin.ID(), Name: plugin.Name(), Description: plugin.Description(), Enabled: plugin.IsEnabled()}
	}
	return plugins
}
