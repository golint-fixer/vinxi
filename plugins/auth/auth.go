package auth

import (
	"errors"
	"net/http"

	"gopkg.in/vinxi/auth.v0"
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/plugin"
)

const (
	// Name defines the plugin semantic identifier.
	Name = "auth"
	// Description defines the plugin friendly description.
	Description = "Authorization and authentication protection"
)

func validator(value interface{}, opts config.Config) error {
	uri := value.(string)
	if uri == "" {
		return errors.New("auth: token cannot be empty")
	}
	return nil
}

// params defines the rule specific configuration params.
var params = plugin.Params{
	plugin.Field{
		Name:        "token",
		Type:        "string",
		Description: "Authentication token",
		Mandatory:   true,
		Validator:   validator,
	},
	plugin.Field{
		Name:        "scheme",
		Type:        "string",
		Description: "Authentication scheme",
		Default:     "Bearer",
		Examples:    []string{"Bearer", "Basic"},
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

// NewToken creates a new auth plugin.
func NewToken(token string) (plugin.Plugin, error) {
	return plugin.NewWithConfig(Plugin, config.Config{"token": token})
}

func handler(opts config.Config) plugin.Handler {
	tokens := []auth.Token{
		{Type: opts.GetString("scheme"), Value: opts.GetString("token")},
	}
	mw := auth.New(&auth.Config{Tokens: tokens})

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw.HandleHTTP(w, r, h)
		})
	}
}

func init() {
	plugin.Register(Plugin)
}
