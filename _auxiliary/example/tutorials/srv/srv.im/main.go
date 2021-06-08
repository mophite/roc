package main

import (
	"fmt"
	"roc/_auxiliary/example/tutorials/proto/pbim"
	"roc/_auxiliary/example/tutorials/srv/srv.im/im"

	"roc/server"
)

func main() {
	var s = server.NewRocServer(server.Namespace("srv.im"))
	pbim.RegisterImServer(s, &im.Im{H: im.NewHub()})
	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
}
