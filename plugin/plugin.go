package plugin

import (
	"net/http"

	"github.com/dchest/uniuri"
	"gopkg.in/vinxi/vinxi.v0/config"
)

// Handler represents the plugin specific HTTP handler function interface.
type Handler func(http.Handler) http.Handler

// Factory functions represents the plugin factory function interface.
type Factory func(config.Config) Plugin

// Plugin represents the required interface implemented by plugins.
type Plugin interface {
	// ID is used to retrieve the plugin unique identifier.
	ID() string
	// Name is used to retrieve the plugin name identifier.
	Name() string
	// Description is used to retrieve a human friendly
	// description of what the plugin does.
	Description() string
	// Config is used to retrieve the rule config.
	Config() config.Config
	// HandleHTTP is used to run the plugin task.
	// Note: add error reporting layer.
	HandleHTTP(http.Handler) http.Handler
}

type plugin struct {
	id          string
	name        string
	description string
	handler     Handler
	config      config.Config
}

// New creates a new Plugin capable interface based on the
// given HTTP handler logic encapsulated as plugin.
func New(name, description string, handler Handler) Factory {
	p := &plugin{id: uniuri.New(), name: name, description: description, handler: handler}
	return func(opts config.Config) Plugin {
		p.config = opts
		return p
	}
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

// Config returns the plugin human readable description about
// what the plugin does and for what it's designed.
func (p *plugin) Config() config.Config {
	return p.config
}

// HandleHTTP implements the required plugin HTTP handler interface
// triggered by the plugin layer during the incoming request call chain.
func (p *plugin) HandleHTTP(h http.Handler) http.Handler {
	next := p.handler(h)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
