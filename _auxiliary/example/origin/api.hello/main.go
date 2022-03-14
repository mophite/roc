package main

import (
    "github.com/go-roc/roc"
    "github.com/go-roc/roc/_auxiliary/example/origin/api.hello/hello"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
)

func main() {
    s := roc.New(
        roc.Namespace("api.hello"),
        roc.HttpApiAddr("0.0.0.0:9999"),
        roc.TCPApiSrvPort(8888),
        roc.WssApiAddr("0.0.0.0:7777","/test/wss"),
    )

    phello.RegisterHelloSrvServer(s.Server(), &hello.Hello{})

    s.Run()
}
