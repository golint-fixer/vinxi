package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err.Error())
	}
}
