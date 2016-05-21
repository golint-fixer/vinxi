package manager

import (
	"net/http"

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
	// Enable is used to enable the current plugin.
	// If the plugin has been already enabled, the call is no-op.
	Enable()
	// Disable is used to disable the current plugin.
	Disable()
	// Remove is used to disable and remove a plugin.
	// Remove()
	// IsEnabled is used to check if a plugin is enabled or not.
	IsEnabled() bool
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
}

// NewPlugin creates a new Plugin capable interface based on the
// given HTTP handler logic encapsulated as plugin.
func NewPlugin(name, description string, handler Handler) Plugin {
	return &plugin{id: uniuri.New(), name: name, description: description, handler: handler}
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

// Disable disables the current plugin.
func (p *plugin) Disable() {
	p.disabled = true
}

// Enable enables the current plugin.
func (p *plugin) Enable() {
	p.disabled = false
}

// IsEnabled returns true if the plugin is enabled.
func (p *plugin) IsEnabled() bool {
	return p.disabled == false
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
	pool []Plugin
}

// NewPluginLayer creates a new plugins layer.
func NewPluginLayer() *PluginLayer {
	return &PluginLayer{}
}

// Use registers one or multiples plugins in the current plugin layer.
func (l *PluginLayer) Use(plugin ...Plugin) {
	l.pool = append(l.pool, plugin...)
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

// Run triggers the plugins layer call chain.
// This function is designed to be executed by top-level middleware layers.
func (l *PluginLayer) Run(w http.ResponseWriter, r *http.Request, h http.Handler) {
	next := h
	for _, plugin := range l.pool {
		next = plugin.HandleHTTP(next)
	}
	next.ServeHTTP(w, r)
}
