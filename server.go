package vinxi

import (
	"net/http"
	"strconv"
	"time"
)

var (
	// DefaultPort stores the default TCP port to listen.
	DefaultPort = 8080

	// DefaultReadTimeout defines the maximum timeout for request read.
	DefaultReadTimeout = 60

	// DefaultWriteTimeout defines the maximum timeout for response write.
	DefaultWriteTimeout = 60
)

// ServerOptions represents the supported server options.
type ServerOptions struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Host         string
	Forward      string
	CertFile     string
	KeyFile      string
}

// Server represents a simple wrapper around http.Server for better convenience
// and easy set up using Vinxi.
type Server struct {
	// Vinxi stores the Vinxi layer instance.
	Vinxi *Vinxi

	// Server stores the http.Server instance.
	Server *http.Server

	// Options stores the server start options.
	Options ServerOptions
}

// NewServer creates a new standard HTTP server.
func NewServer(o ServerOptions) *Server {
	// Apply default options
	if o.Port == 0 {
		o.Port = DefaultPort
	}
	if o.ReadTimeout == 0 {
		o.ReadTimeout = DefaultReadTimeout
	}
	if o.WriteTimeout == 0 {
		o.WriteTimeout = DefaultWriteTimeout
	}

	addr := o.Host + ":" + strconv.Itoa(o.Port)
	svr := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(o.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(o.WriteTimeout) * time.Second,
	}

	vinxi := New()
	vinxi.BindServer(svr)

	if o.Forward != "" {
		vinxi.Forward(o.Forward)
	}

	return &Server{
		Options: o,
		Server:  svr,
		Vinxi:   vinxi,
	}
}

// Forward defines the default URL to forward incoming traffic.
func (s *Server) Forward(uri string) *Server {
	s.Vinxi.Forward(uri)
	return s
}

// Use attaches a new middleware handler for incoming HTTP traffic.
func (s *Server) Use(handler interface{}) *Server {
	s.Vinxi.Use(handler)
	return s
}

// UsePhase attaches a new middleware handler to a specific phase.
func (s *Server) UsePhase(phase string, handler interface{}) *Server {
	s.Vinxi.UsePhase(phase, handler)
	return s
}

// UseFinalHandler uses a new middleware handler function as final handler.
func (s *Server) UseFinalHandler(fn http.Handler) *Server {
	s.Vinxi.UseFinalHandler(fn)
	return s
}

// Listen starts listening on network.
func (s *Server) Listen() error {
	if s.Options.CertFile != "" && s.Options.KeyFile != "" {
		return s.Server.ListenAndServeTLS(s.Options.CertFile, s.Options.KeyFile)
	}
	return s.Server.ListenAndServe()
}
