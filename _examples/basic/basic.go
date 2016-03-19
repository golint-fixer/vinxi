package main

import (
	"fmt"
	"gopkg.in/vinci-proxy/vinci.v0"
)

func main() {
	fmt.Printf("Server listening on port: %d\n", 3100)
	vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

	vs.Vinci.Forward("http://localhost:8080")

	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
