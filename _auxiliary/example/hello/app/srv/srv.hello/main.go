package main

import (
    "github.com/go-roc/roc/_auxiliary/example/httpservice/app/api/api.hello/hello"
    "github.com/go-roc/roc/_auxiliary/example/httpservice/proto/phello"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service"
)

func main() {
    s := service.New(service.Namespace("srv.hello"))

    phello.RegisterHelloWorldServer(s.Server(), &hello.Hello{})
    err := s.Run()
    if err != nil {
        rlog.Error(err)
    }
}
