package main

import (
	"fmt"
	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/intercept"
	"strings"
	// "io/ioutil"
	// "net/http"
)

func main() {
	fmt.Printf("Server listening on port: %d\n", 3100)
	vs := vinxi.NewServer(vinxi.ServerOptions{Host: "localhost", Port: 3100})

	vs.Vinci.Use(intercept.Request(func(req *intercept.RequestModifier) {
		str, _ := req.ReadString()
		fmt.Printf("Request body: %s \n", str)
		req.String("foo bar")
	}))

	vs.Vinci.Use(intercept.Response(func(res *intercept.ResponseModifier) {
		data, _ := res.ReadString()
		fmt.Printf("Response body: %s \n", data)
		str := strings.Replace(data, "The MIT License", "Apache License", 1)
		res.String(str)
	}))

	vs.Vinci.Forward("http://localhost:8080")

	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
