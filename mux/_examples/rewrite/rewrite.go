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
	m.If(mux.MatchMethod("GET"), mux.MatchPath("^/ip"))

	m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		// Overwrite request URI
		r.RequestURI = "/get"
		h.ServeHTTP(w, r)
	})

	vs.Use(m)
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", 3100)
	err := vs.Listen()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
