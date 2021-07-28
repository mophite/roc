package main

import (
    "github.com/go-roc/roc/_auxiliary/example/hello/app/api/api.hello/hello"
    "github.com/go-roc/roc/_auxiliary/example/hello/proto/phello"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service"
)

func main() {
    s := service.New(
        service.HttpAddress("0.0.0.0:9999"),
        service.Namespace("api.hello"),
        service.TCPAddress("0.0.0.0:8888"),
        service.WssAddress("0.0.0.0:10000", "/hello"),
    )

    phello.RegisterHelloWorldServer(s.Server(), &hello.Hello{Client: s.Client()})
    err := s.Run()
    if err != nil {
        rlog.Error(err)
    }
}
