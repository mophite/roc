package main

import (
	"fmt"
	"roc/_auxiliary/example/tutorials/proto/pbhello"
	"roc/_auxiliary/example/tutorials/srv/srv.hello/hello"
	"roc/internal/registry"
	"roc/server"
)

func main() {
	var opt = registry.Address([]string{"82.157.14.79:2379"})

	var s = server.NewRocServer(
		server.Namespace("srv.hello"),
		server.TCPAddress("127.0.0.1:8089"),
		server.Registry(registry.NewRegistry(opt)),
	)
	pbhello.RegisterHelloWorldServer(s, &hello.Hello{})
	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
}
