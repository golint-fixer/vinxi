package main

import (
	"fmt"
	"gopkg.in/vinci-proxy/vinci.v0"
)

func main() {
	server := vinci.New(vinci.ServerOptions{Host: "localhost", Post: 3100})

	err := server.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
