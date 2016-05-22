package manager

import (
	"net/http"
	"sync"

	"github.com/dchest/uniuri"
	"gopkg.in/vinxi/vinxi.v0/layer"
)

// Handler represents the plugin specific HTTP handler function interface.
type Handler func(http.Handler) http.Handler

// Plugin represents the required interface implemented by plugins.
type Plugin interface {
	// ID is used to retrieve the plugin unique identifier.
	ID() string
	// Name is used to retrieve the plugin name identifier.
	Name() string
	// Description is used to retrieve a human friendly
	// description of what the plugin does.
	Description() string
	// JSONConfig is used to retrieve the plugin specific
	// config as serialized JSON notation.
	JSONConfig() string
	// HandleHTTP is used to run the plugin task.
	// Note: add erro reporting layer
	HandleHTTP(http.Handler) http.Handler
}

type plugin struct {
	disabled    bool
	id          string
	name        string
	description string
	handler     Handler
	config      interface{}
}

// NewPlugin creates a new Plugin capable interface based on the
// given HTTP handler logic encapsulated as plugin.
func NewPlugin(name, description string, handler Handler) Plugin {
	return &plugin{id: uniuri.New(), name: name, description: description, handler: handler}
}

func (p *plugin) Configure(config interface{}) {
	p.config = config
}

// ID returns the plugin identifer.
func (p *plugin) ID() string {
	return p.id
}

// Name returns the plugin semantic name identifier.
func (p *plugin) Name() string {
	return p.name
}

// Description returns the plugin human readable description about
// what the plugin does and for what it's designed.
func (p *plugin) Description() string {
	return p.description
}

// JSONConfig returns the plugin human readable description about
// what the plugin does and for what it's designed.
func (p *plugin) JSONConfig() string {
	return p.description
}

// HandleHTTP implements the required plugin HTTP handler interface
// triggered by the plugin layer during the incoming request call chain.
func (p *plugin) HandleHTTP(h http.Handler) http.Handler {
	next := p.handler(h)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// PluginLayer represents a plugins layer designed to intrument
// proxies providing plugin based dynamic configuration
// capabilities, such as register/unregister or
// enable/disable plugins at runtime satefy.
type PluginLayer struct {
	rwm  sync.RWMutex
	pool []Plugin
}

// NewPluginLayer creates a new plugins layer.
func NewPluginLayer() *PluginLayer {
	return &PluginLayer{}
}

// Use registers one or multiples plugins in the current plugin layer.
func (l *PluginLayer) Use(plugin ...Plugin) {
	l.rwm.Lock()
	l.pool = append(l.pool, plugin...)
	l.rwm.Unlock()
}

// Len returns the registered plugins length.
func (l *PluginLayer) Len() int {
	return len(l.pool)
}

// Register implements the middleware Register method.
func (l *PluginLayer) Register(mw *layer.Layer) {
	mw.Use("error", l.Run)
	mw.Use("request", l.Run)
}

// Remove removes a plugin looking by its unique identifier.
func (l *PluginLayer) Remove(id string) bool {
	l.rwm.Lock()
	defer l.rwm.Unlock()

	for i, plugin := range l.pool {
		if plugin.ID() == id {
			l.pool = append(l.pool[:i], l.pool[i+1:]...)
			return true
		}
	}

	return false
}

// Run triggers the plugins layer call chain.
// This function is designed to be executed by top-level middleware layers.
func (l *PluginLayer) Run(w http.ResponseWriter, r *http.Request, h http.Handler) {
	next := h
	l.rwm.RLock()
	for _, plugin := range l.pool {
		next = plugin.HandleHTTP(next)
	}
	l.rwm.RUnlock()
	next.ServeHTTP(w, r)
}
