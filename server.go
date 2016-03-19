package vinci

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
	CertFile     string
	KeyFile      string
}

// Server represents a simple wrapper around http.Server for better convenience
// and easy set up using Vinci.
type Server struct {
	// Vinci stores the Vinci layer instance.
	Vinci *Vinci

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

	vinci := New()
	vinci.BindServer(svr)

	return &Server{
		Options: o,
		Server:  svr,
		Vinci:   vinci,
	}
}

// Forward defines the default URL to forward incoming traffic.
func (s *Server) Forward(uri string) *Server {
	s.Vinci.Forward(uri)
	return s
}

// Use attaches a new middleware handler for incoming HTTP traffic.
func (s *Server) Use(handler interface{}) *Server {
	s.Vinci.Middleware.Use(handler)
	return s
}

// UseError attaches a new middleware handler to the.
func (s *Server) UseError(handler interface{}) *Server {
	s.Vinci.Middleware.UseError(handler)
	return s
}

// UsePhase attaches a new middleware handler to a specific phase.
func (s *Server) UsePhase(phase string, handler interface{}) *Server {
	s.Vinci.Middleware.UsePhase(phase, handler)
	return s
}

// UseFinalHandler uses a new middleware handler function as final handler.
func (s *Server) UseFinalHandler(fn http.Handler) *Server {
	s.Vinci.Middleware.UseFinalHandler(fn)
	return s
}

// Listen starts listening on network.
func (s *Server) Listen() error {
	if s.Options.CertFile != "" && s.Options.KeyFile != "" {
		return s.Server.ListenAndServeTLS(s.Options.CertFile, s.Options.KeyFile)
	}
	return s.Server.ListenAndServe()
}
