package static

import (
	"errors"
	"net/http"
	"os"

	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/plugin"
)

const (
	// Name defines the plugin semantic identifier.
	Name = "static"
	// Description defines the plugin friendly description.
	Description = "serve local directory"
)

func isDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	return file.IsDir()
}

// params defines the rule specific configuration params.
var params = plugin.Params{
	plugin.Field{
		Name:        "path",
		Type:        "string",
		Description: "Local path to serve",
		Mandatory:   true,
		Validator: func(value interface{}, opts config.Config) error {
			path := value.(string)
			if path == "" {
				return errors.New("static: path cannot be empty")
			}
			if !isDir(path) {
				return errors.New("static: path does not exists or invalid")
			}
			return nil
		},
	},
}

// Plugin exposes the rule metadata information.
// Mostly used internally.
var Plugin = plugin.Info{
	Name:        Name,
	Description: Description,
	Factory:     factory,
	Params:      params,
}

// factory represents the rule factory function
// designed to be called via rules constructor.
func factory(opts config.Config) plugin.Handler {
	return handler(opts)
}

// New creates a new static plugin who serves
// files of the given server local path.
func New(path string) (plugin.Plugin, error) {
	return plugin.NewWithConfig(Plugin, config.Config{"path": path})
}

func handler(opts config.Config) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.FileServer(http.Dir(opts.GetString("path")))
	}
}

func init() {
	plugin.Register(Plugin)
}
