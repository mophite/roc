package main

import (
	"github.com/go-roc/roc/_auxiliary/example/hello/app/srv/srv.hello/hello"
	"github.com/go-roc/roc/_auxiliary/example/hello/proto/phello"
	"github.com/go-roc/roc/rlog"
	"github.com/go-roc/roc/service"
)

func main() {
	s := service.New(
		service.TCPAddress("0.0.0.0:11111"),
		service.Namespace("srv.hello"),
	)

	phello.RegisterHelloWorldServer(s.Server(), &hello.Hello{})
	err := s.Run()
	if err != nil {
		rlog.Error(err)
	}
}
