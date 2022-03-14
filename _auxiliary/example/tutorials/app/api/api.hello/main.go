package main

import (
    "github.com/go-roc/roc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/app/api/api.hello/hello"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/app/api/api.hello/say"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/app/api/api.hello/upload"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/internal/ipc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
)

func main() {
    s := roc.New(
        roc.RootRouterRedirect("/hello/hello/say"),
        roc.Namespace("api.hello"),
        roc.HttpApiAddr("0.0.0.0:10000"),
        roc.TCPApiSrvPort(10001),
        roc.WssApiAddr("0.0.0.0:10002", "/hello"),
    )

    phello.RegisterHelloServer(s.Server(), &say.Say{})
    phello.RegisterFileServer(s.Server(), &upload.File{})
    phello.RegisterHelloSrvServer(s.Server(), &hello.Hello{})

    ipc.InitIpc(s)

    //imChannel.Channel()
    //imStrem.Stream()
    //ImTest.Im()
    s.Run()
}
