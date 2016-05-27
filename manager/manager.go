package manager

import (
	"net/http"
	"sync"

	"github.com/h2non/httprouter"
	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/layer"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"

	// An empty import is required to load all the rules subpackages
	_ "gopkg.in/vinxi/vinxi.v0/plugins"
	_ "gopkg.in/vinxi/vinxi.v0/rules"
)

// Manager represents the vinxi proxy admin manager.
type Manager struct {
	sm        sync.RWMutex
	scopes    []*Scope
	im        sync.RWMutex
	instances []*Instance

	// AdminPlugins stores the HTTP admin server plugins.
	AdminPlugins *plugin.Layer
	// Plugins stores the global plugin layer.
	Plugins *plugin.Layer
	// Server stores the HTTP server used for the admin.
	Server *http.Server
	// Layer stores the manager internal middleware layer.
	Layer *layer.Layer
	// Router stores the manager HTTP router for the admin server.
	Router *httprouter.Router
}

// New creates a new manager able to manage
// and configure multiple vinxi proxy instance.
func New() *Manager {
	mw := layer.New()
	aplugins := plugin.NewLayer()

	// Define middleware priority
	mw.UsePriority(aplugins, layer.TailPriority)

	return &Manager{
		Layer:        mw,
		AdminPlugins: aplugins,
		Router:       httprouter.New(),
		Plugins:      plugin.NewLayer(),
	}
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
func (m *Manager) Manage(name, description string, proxy *vinxi.Vinxi) *Instance {
	// Register the managed Vinxi instance
	instance := NewInstance(name, description, proxy)

	// Register manager middleware in the proxy
	proxy.Layer.UsePriority(layer.RequestPhase, layer.Tail, m)

	// Register the vinxi instance specific middleware layer
	proxy.Layer.UsePriority(layer.RequestPhase, layer.Tail, instance)

	// Register instance
	m.im.Lock()
	m.instances = append(m.instances, instance)
	m.im.Unlock()

	return instance
}

// Use attaches a new middleware handler for incoming HTTP traffic.
func (m *Manager) Use(handler ...interface{}) {
	m.Layer.Use(layer.RequestPhase, handler...)
}

// UsePhase attaches a new middleware handler to a specific phase.
func (m *Manager) UsePhase(phase string, handler ...interface{}) {
	m.Layer.Use(phase, handler...)
}

// UseFinalHandler uses a new middleware handler function as final handler.
func (m *Manager) UseFinalHandler(fn http.Handler) {
	m.Layer.UseFinalHandler(fn)
}

// UseAdminPlugin registers one or multiple plugins at manager admin level.
func (m *Manager) UseAdminPlugin(plugins ...plugin.Plugin) {
	m.AdminPlugins.Use(plugins...)
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

// UseScope registers one or multiple scopes at global manager level.
func (m *Manager) UseScope(scopes ...*Scope) {
	m.sm.Lock()
	m.scopes = append(m.scopes, scopes...)
	m.sm.Unlock()
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
	scope := m.NewScope("default", "Default scope")
	scope.UseRule(rules...)
	return scope
}

// UsePlugin registers one or multiple plugins at global manager level.
func (m *Manager) UsePlugin(plugins ...plugin.Plugin) {
	m.Plugins.Use(plugins...)
}

// GetPlugin finds and returns a plugin by its ID or name.
func (m *Manager) GetPlugin(name string) plugin.Plugin {
	return m.Plugins.Get(name)
}

// RemovePlugin removes a plugin by its ID.
func (m *Manager) RemovePlugin(id string) bool {
	return m.Plugins.Remove(id)
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

// Scopes returns the registered scopes at global level.
func (m *Manager) Scopes() []*Scope {
	m.sm.Lock()
	defer m.sm.Unlock()
	return m.scopes
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

// Instances returns the registered vinxi instances in the manager.
func (m *Manager) Instances() []*Instance {
	m.im.Lock()
	defer m.im.Unlock()
	return m.instances
}

// GetInstance finds and returns a vinxi managed instance.
func (m *Manager) GetInstance(name string) *Instance {
	m.im.Lock()
	defer m.im.Unlock()

	for _, instance := range m.instances {
		if instance.ID() == name || instance.Metadata().Name == name {
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
		if instance.ID() == name || instance.Metadata().Name == name {
			m.instances = append(m.instances[:i], m.instances[i+1:]...)
			return true
		}
	}

	return false
}

// HandleHTTP is triggered by the vinxi middleware layer on incoming HTTP request.
func (m *Manager) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	// Declare the final handler in the call chain
	next := h

	// Build the scope handlers call chain
	m.sm.RLock()
	for _, scope := range m.scopes {
		next = http.HandlerFunc(scope.HandleHTTP(next))
	}
	m.sm.RUnlock()

	// Run global plugins, then global scopes
	m.Plugins.Run(w, r, next)
}

// serveHTTP is used to handle HTTP traffic via admin API.
func (m *Manager) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// trigger middleware layer first, then the router
	m.Layer.Run(layer.RequestPhase, w, r, m.Router)
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
