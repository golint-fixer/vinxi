package main

import (
	"gopkg.in/vinxi/vinxi.v0"
)

func main() {
	// Creates a new vinxi proxy
	v := vinxi.New()
	v.Forward("http://localhost:9090")
	v.ListenAndServe(vinxi.ServerOptions{})
}
