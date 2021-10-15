package main

import (
    "github.com/go-roc/roc/_auxiliary/example/origin/api.hello/hello"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/service"
)

func main() {
    s := service.New(
        service.Namespace("api.hello"),
        service.HttpAddress("0.0.0.0:9999"),
        service.TCPAddress("0.0.0.0:8888"),
    )

    phello.RegisterHelloSrvServer(s.Server(), &hello.Hello{})

    s.Run()
}
