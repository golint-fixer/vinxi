package forward

import (
	"errors"
	"net/http"
	"net/url"

	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/forward"
	"gopkg.in/vinxi/vinxi.v0/plugin"
)

const (
	// Name defines the plugin semantic identifier.
	Name = "forward"
	// Description defines the plugin friendly description.
	Description = "Forward HTTP traffic to remote servers"
)

func validator(value interface{}, opts config.Config) error {
	uri := value.(string)
	if uri == "" {
		return errors.New("forward: url param cannot be empty")
	}
	_, err := url.Parse(uri)
	if err != nil {
		return errors.New("forward: invalid URL (" + err.Error() + ")")
	}
	return nil
}

// params defines the rule specific configuration params.
var params = plugin.Params{
	plugin.Field{
		Name:        "url",
		Type:        "string",
		Description: "Target server URL",
		Mandatory:   true,
		Validator:   validator,
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

// New creates a new forward plugin.
func New(url string) (plugin.Plugin, error) {
	return plugin.NewWithConfig(Plugin, config.Config{"url": url})
}

func handler(opts config.Config) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(forward.To(opts.GetString("url")))
	}
}

func init() {
	plugin.Register(Plugin)
}
