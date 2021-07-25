package main

import (
	"github.com/go-roc/roc/_auxiliary/example/hello/app/api/api.hello/hello"
	"github.com/go-roc/roc/_auxiliary/example/hello/proto/phello"
	"github.com/go-roc/roc/service"
)

func main() {
	s := service.New(
		service.TCPAddress("0.0.0.0:8888"),
		service.HttpAddress("0.0.0.0:9999"),
	)

	phello.RegisterHelloWorldServer(s.Server(), &hello.Hello{})
	err := s.Run()
	if err != nil {
		panic(err)
	}
}
