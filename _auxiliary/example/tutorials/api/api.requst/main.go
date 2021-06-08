package main

import (
	"net/http"
	"roc/_auxiliary/example/tutorials/api/api.requst/hello"
)

// open with browser:
// http://localhost:8080/say
func main() {
	h := hello.NewHello()
	http.HandleFunc("/say", h.SayHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
