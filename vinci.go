package vinci

import (
	"gopkg.in/vinci-proxy/forward.v0"
	"gopkg.in/vinci-proxy/layer.v0"
	"net/http"
)

// DefaultForwarder stores the default http.Handler to be used to forward the traffic.
// By default the proxy will reply with 502 Bad Gateway if no custom forwarder is defined.
var DefaultForwarder = layer.FinalHandler

// Vinci represents the vinci proxy structure.
type Vinci struct {
	// Layer stores the proxy top-level middleware layer.
	Layer *layer.Layer
}

// New creates a new vinci proxy layer.
func New() *Vinci {
	return &Vinci{Layer: layer.New()}
}

// Forward defines the default URL to forward incoming traffic.
func (v *Vinci) Forward(uri string) *Vinci {
	return v.UseForwarder(http.HandlerFunc(forward.To(uri)))
}

// UseForwarder uses a custom forwarder HTTP handler to proxy incoming traffic.
func (v *Vinci) UseForwarder(forwarder http.Handler) *Vinci {
	v.Layer.UseFinalHandler(forwarder)
	return v
}

// Use attaches a new middleware handler for incoming HTTP traffic.
func (v *Vinci) Use(handler interface{}) *Vinci {
	v.Layer.Use("request", handler)
	return v
}

// UsePhase attaches a new middleware handler to a specific phase.
func (v *Vinci) UsePhase(phase string, handler interface{}) *Vinci {
	v.Layer.Use(phase, handler)
	return v
}

// UseFinalHandler uses a new middleware handler function as final handler.
func (v *Vinci) UseFinalHandler(fn http.Handler) *Vinci {
	v.Layer.UseFinalHandler(fn)
	return v
}

// Flush flushes all the middleware stack.
func (v *Vinci) Flush() *Vinci {
	v.Layer.Flush()
	return v
}

// BindServer binds the vinci handler to the given http.Server.
func (v *Vinci) BindServer(server *http.Server) *Vinci {
	server.Handler = v
	return v
}

// ServeHTTP implements the required http.Handler interface.
func (v *Vinci) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	v.Layer.Run("request", w, req, nil)
}
