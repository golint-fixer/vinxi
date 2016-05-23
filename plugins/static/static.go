package static

import (
	"net/http"

	"gopkg.in/vinxi/vinxi.v0/plugin"
)

// New creates a new static plugin who serves
// files of the given server local path.
func New(path string) plugin.Plugin {
	return plugin.New("static", "serve static files", staticHandler(path))(map[string]interface{}{"path": path})
}

func staticHandler(path string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.FileServer(http.Dir(path))
	}
}
