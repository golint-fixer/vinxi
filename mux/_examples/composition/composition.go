package main

import (
	"fmt"
	"gopkg.in/vinxi/mux.v0"
	"gopkg.in/vinxi/vinxi.v0"
	"net/http"
)

func main() {
	vs := vinxi.NewServer(vinxi.ServerOptions{Host: "localhost", Port: 3100})

	// Create a custom multiplexer for /ip path
	ip := mux.If(mux.Path("^/ip"))
	ip.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	// Create a custom multiplexer for /headers path
	headers := mux.If(mux.Path("^/headers"))
	headers.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Errorf("Headers: %#v", r.Header).Error()))
	})

	// Creates the root multiplexer who host both multiplexers
	m := mux.New()
	m.If(mux.MatchMethod("GET"))
	m.Use(ip)
	m.Use(headers)

	// Register the multiplexer in the vinxi
	vs.Use(m)
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", 3100)
	err := vs.Listen()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
