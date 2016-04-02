package main

import (
	"fmt"
	"gopkg.in/vinxi/vinxi.v0"
)

func main() {
	server := vinxi.New(vinxi.ServerOptions{Host: "localhost", Post: 3100})

	err := server.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
