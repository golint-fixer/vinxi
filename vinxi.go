package vinxi

import (
	"gopkg.in/vinxi/forward.v0"
	"gopkg.in/vinxi/layer.v0"
	"gopkg.in/vinxi/router.v0"
	"net/http"
)

// DefaultForwarder stores the default http.Handler to be used to forward the traffic.
// By default the proxy will reply with 502 Bad Gateway if no custom forwarder is defined.
var DefaultForwarder = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fwd, _ := forward.New(forward.PassHostHeader(true))
	fwd.ServeHTTP(w, r)
})

// Middleware defines the required interface implemented
// by public middleware capable entities in the vinxi ecosystem.
type Middleware interface {
	// Use is used to register one or multiple middleware handlers.
	Use(...interface{}) Middleware

	// UsePhase is used to register one or multiple middleware
	// handlers for a specific middleware phase.
	UsePhase(string, ...interface{}) Middleware

	// UseFinalHandler is used to register the final request handler
	// usually to define the error or forward handlers.
	UseFinalHandler(http.Handler) Middleware
}

// Route represents the required route capable interface
type Route interface {
	Middleware
	http.Handler
	Forward(string) Route
	Handle(http.HandlerFunc)
}

// Router represents the router capable interface.
type Router interface {
	Route(string, string) Route
	Match(string string) (Route, error)
}

// Vinxi represents the vinxi proxy structure.
type Vinxi struct {
	// Layer stores the proxy top-level middleware layer.
	Layer *layer.Layer

	// Router stores the built-in router.
	Router *router.Router
}

// New creates a new vinxi proxy layer.
func New() *Vinxi {
	v := &Vinxi{Layer: layer.New(), Router: router.New()}
	v.Layer.UsePriority("request", layer.Tail, v.Router)
	v.UseFinalHandler(DefaultForwarder)
	return v
}

// Get will register a pattern for GET requests.
// It also registers pat for HEAD requests. If this needs to be overridden, use
// Head before Get with pat.
func (v *Vinxi) Get(path string) *router.Route {
	return r.Route("GET", path)
}

// Post will register a pattern for POST requests.
func (v *Vinxi) Post(path string) *router.Route {
	return r.Route("POST", path)
}

// Put will register a pattern for PUT requests.
func (v *Vinxi) Put(path string) *router.Route {
	return r.Route("PUT", path)
}

// Delete will register a pattern for DELETE requests.
func (v *Vinxi) Delete(path string) *router.Route {
	return r.Route("DELETE", path)
}

// Options will register a pattern for OPTIONS requests.
func (v *Vinxi) Options(path string) *router.Route {
	return r.Route("OPTIONS", path)
}

// Patch will register a pattern for PATCH requests.
func (v *Vinxi) Patch(path string) *router.Route {
	return r.Route("PATCH", path)
}

// All will register a pattern for any HTTP method.
func (v *Vinxi) All(path string) *router.Route {
	return r.Route("*", path)
}

// Route will register a new route for the given pattern and HTTP method.
func (v *Vinxi) Route(method, path string) *router.Route {
	return r.Router.Route(method, path)
}

// Forward defines the default URL to forward incoming traffic.
func (v *Vinxi) Forward(uri string) *Vinxi {
	return v.UseFinalHandler(http.HandlerFunc(forward.To(uri)))
}

// Use attaches a new middleware handler for incoming HTTP traffic.
func (v *Vinxi) Use(handler ...interface{}) *Vinxi {
	v.Layer.Use(layer.RequestPhase, handler...)
	return v
}

// UsePhase attaches a new middleware handler to a specific phase.
func (v *Vinxi) UsePhase(phase string, handler ...interface{}) *Vinxi {
	v.Layer.Use(phase, handler...)
	return v
}

// UseFinalHandler uses a new middleware handler function as final handler.
func (v *Vinxi) UseFinalHandler(fn http.Handler) *Vinxi {
	v.Layer.UseFinalHandler(fn)
	return v
}

// Flush flushes all the middleware stack.
func (v *Vinxi) Flush() {
	v.Layer.Flush()
}

// BindServer binds the vinxi HTTP handler to the given http.Server.
func (v *Vinxi) BindServer(server *http.Server) {
	server.Handler = v
}

// ServeHTTP implements the required http.Handler interface.
func (v *Vinxi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	v.Layer.Run("request", w, req, nil)
}
