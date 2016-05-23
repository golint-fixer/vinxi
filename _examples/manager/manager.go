package main

import (
	"fmt"
	"net/http"

	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/manager"
	"gopkg.in/vinxi/vinxi.v0/plugins/static"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

const port = 3100

func main() {
	// Creates a new vinxi proxy
	v := vinxi.New()

	// Creates a new manager for the vinxi proxy
	mgr := manager.New()
	mgr.Manage("default", "This a default server instance", v)

	// Starts default admin HTTP server
	go mgr.ServeDefault()

	// Register scopes
	scope := mgr.NewScope("default", "Default scope")
	scope.UseRule(rule.Init("path", map[string]interface{}{"path": "/(.*)"}))
	scope.UseRule(rule.Init("vhost", map[string]interface{}{"host": "localhost"}))
	scope.UsePlugin(static.New("/Users/h2non/Projects/vinxi"))

	// Registers a simple middleware handler
	v.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("Server", "vinxi")
		h.ServeHTTP(w, r)
	})

	// Forward traffic to httpbin.org by default
	v.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", port)
	_, err := v.ListenAndServe(vinxi.ServerOptions{Port: port})
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
