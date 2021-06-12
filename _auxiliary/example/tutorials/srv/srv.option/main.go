package main

import (
	"fmt"
	"roc/_auxiliary/example/tutorials/proto/pbhello"
	"roc/_auxiliary/example/tutorials/srv/srv.hello/hello"

	"roc/server"
)

func main() {
	var s = server.NewRocServer(server.Namespace("srv.hello"))
	pbhello.RegisterHelloWorldServer(s, &hello.Hello{})
	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
}
