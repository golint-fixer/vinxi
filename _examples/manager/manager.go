package main

import (
	"fmt"
	"net/http"

	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/manager"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

const port = 3100

func main() {
	// Creates a new manager for the vinxi proxy
	mgr := manager.New()

	// Attach global plugin
	// plu, err := plugin.Init("auth", config.Config{"token": "foo"})
	// if err != nil {
	// fmt.Printf("Error: %s\n", err)
	// return
	// }
	// mgr.UsePlugin(plu)

	// Creates a new vinxi proxy and manage it
	v := vinxi.New()
	instance := mgr.Manage("default", "This a default proxy", v)

	// Instance level scope
	scope := instance.NewScope("custom", "Custom scope")
	scope.UseRule(rule.Init("path", config.Config{"path": "/image/(.*)"}))

	plu, err := plugin.Init("forward", config.Config{"url": "http://httpbin.org"})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	scope.UsePlugin(plu)

	// Starts default admin HTTP server
	go mgr.ServeDefault()

	// Register scopes
	scope = mgr.NewScope("default", "Default scope")
	scope.UseRule(rule.Init("path", config.Config{"path": "/vinxi/(.*)"}))
	scope.UseRule(rule.Init("vhost", config.Config{"host": "localhost"}))

	plu, err = plugin.Init("static", config.Config{"path": "/Users/h2non/Projects/vinxi"})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	scope.UsePlugin(plu)

	// Registers a simple middleware handler
	v.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("Server", "vinxi")
		h.ServeHTTP(w, r)
	})

	// Forward traffic to httpbin.org by default
	v.Forward("http://www.apache.org")

	fmt.Printf("Server listening on port: %d\n", port)
	_, err = v.ListenAndServe(vinxi.ServerOptions{Port: port})
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
