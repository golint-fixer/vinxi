package main

import (
	"fmt"
	"gopkg.in/vinci-proxy/vinci.v0"
)

func main() {
	mux := vinci.New(vinci.ServerOptions{Host: "localhost"})

	mux.Route("/").
		Forward("http://foo.com").
		Use(middleware.SetHeader())

	err := server.Listen(3100)
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
