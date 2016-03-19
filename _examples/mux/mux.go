package main

import (
	"fmt"
	"gopkg.in/vinci-proxy/mux.v0"
	"gopkg.in/vinci-proxy/vinci.v0"
	"gopkg.in/vinci-proxy/vinci.v0/middleware/forward"
	"gopkg.in/vinci-proxy/vinci.v0/route"
	"net/http"
)

func main() {
	fmt.Printf("Server listening on port: %d\n", 3100)
	vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

	m := mux.New()
	m.If(mux.MatchHost("localhost:3100"))

	m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("Server", "vinci")
		h.ServeHTTP(w, r)
	})

	m.Use(route.New("/foo").Handler(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("foo bar"))
	}))

	m.Use(forward.To("http://127.0.0.1:8080"))

	vs.Vinci.Use(m)
	vs.Vinci.Forward("http://127.0.0.1")

	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
