package plugin

import (
	"net/http"

	"github.com/dchest/uniuri"
	"gopkg.in/vinxi/vinxi.v0/config"
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
	// Config is used to retrieve the user defined plugin config.
	Config() config.Config
	// Metadata is used to retrieve the plugin metadata.
	Metadata() config.Config
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
	metadata    config.Config
}

// New creates a new Plugin capable interface based on the
// given HTTP handler logic encapsulated as plugin.
func New(info Info) FactoryFunc {
	return func(opts config.Config) (Plugin, error) {
		return NewWithConfig(info, opts)
	}
}

// NewWithConfig creates a new Plugin capable interface based on the
// given HTTP handler logic encapsulated as plugin.
func NewWithConfig(info Info, opts config.Config) (Plugin, error) {
	if err := Validate(info.Params, opts); err != nil {
		return nil, err
	}

	return &plugin{
		id:          uniuri.New(),
		name:        info.Name,
		description: info.Description,
		config:      opts,
		handler:     info.Factory(opts),
	}, nil
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

// Config returns the plugin human readable description about
// what the plugin does and for what it's designed.
func (p *plugin) Metadata() config.Config {
	return p.metadata
}

// HandleHTTP implements the required plugin HTTP handler interface
// triggered by the plugin layer during the incoming request call chain.
func (p *plugin) HandleHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.handler(h).ServeHTTP(w, r)
	})
}
