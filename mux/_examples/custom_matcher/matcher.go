package main

import (
	"fmt"
	"gopkg.in/vinxi/mux.v0"
	"gopkg.in/vinxi/vinxi.v0"
	"net/http"
)

func main() {
	vs := vinxi.NewServer(vinxi.ServerOptions{Host: "localhost", Port: 3100})

	m := mux.New()

	// Register a custom matcher function
	m.If(func(req *http.Request) bool {
		return req.Method == "GET" && req.RequestURI == "/foo"
	})

	m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("Server", "vinxi")
		h.ServeHTTP(w, r)
	})

	m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Write([]byte("foo"))
	})

	vs.Use(m)
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", 3100)
	err := vs.Listen()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
