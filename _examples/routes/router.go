package main

import (
	"fmt"
	"gopkg.in/vinci-proxy/vinci.v0"
	"gopkg.in/vinci-proxy/vinci.v0/route"
	"net/http"
)

func main() {
	fmt.Printf("Server listening on port: %d\n", 3100)
	vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

	r := route.New("/")

	r.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("Server", "vinci")
		h.ServeHTTP(w, r)
	})

	r.UseError(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("Server", "vinci")
		w.WriteHeader(500)
		w.Write([]byte("server error"))
	})

	r.Forward("http://localhost:8080")

	vs.Vinci.Use(r)

	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
