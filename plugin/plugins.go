package plugin

import (
	"errors"

	"gopkg.in/vinxi/vinxi.v0/config"
)

// ErrPluginNotFound is used when a plugin does not exists.
var ErrPluginNotFound = errors.New("vinxi: plugin does not exists")

// Plugins is used to store the available plugins globally.
var Plugins = make(map[string]Info)

// Factory represents the plugin factory function interface.
type Factory func(config.Config) Handler

// NewFunc represents the Plugin constructor factory function interface.
type NewFunc func(config.Config) (Plugin, error)

// Validator represents the plugin config field validator function interface.
type Validator func(interface{}, config.Config) error

// Info represents the plugin entity fields
// storing the name, description and factory function
// used to initialize the fields.
type Info struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Params      Params  `json:"params,omitempty"`
	Factory     Factory `json:"-"`
}

// Params represents the list of supported config fields by plugins.
type Params []Field

// Field is used to declare specific config fields supported by plugins.
type Field struct {
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Mandatory   bool        `json:"mandatory,omitempty"`
	Examples    []string    `json:"examples,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Validator   Validator   `json:"-"`
}

// Register registers the given plugin in the current store.
func Register(plugin Info) {
	Plugins[plugin.Name] = plugin
}

// Init is used to initialize a new plugin by name identifier
// based on the given config options.
func Init(name string, opts config.Config) (Plugin, error) {
	if !Exists(name) {
		return nil, ErrPluginNotFound
	}
	return NewWithConfig(Plugins[name], opts)
}

// Get is used to find and retrieve a plugin.
func Get(name string) NewFunc {
	plugin, ok := Plugins[name]
	if ok {
		return New(plugin)
	}
	return nil
}

// GetFactory is used to find and retrieve a plugin factory function.
func GetFactory(name string) Factory {
	plugin, ok := Plugins[name]
	if ok {
		return plugin.Factory
	}
	return nil
}

// GetInfo is used to find and retrieve a plugin info struct, if exists.
func GetInfo(name string) Info {
	return Plugins[name]
}

// Exists is used to check if a given plugin name exists.
func Exists(name string) bool {
	_, ok := Plugins[name]
	return ok
}
