package main

import (
    ImTest "github.com/go-roc/roc/_auxiliary/example/tutorials/app/api/api.hello/im"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/app/api/api.hello/say"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/app/api/api.hello/upload"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/internal/ipc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/service"
)

func main() {
    s := service.New(
        service.Namespace("api.hello"),
        service.HttpAddress("0.0.0.0:9999"),
        service.TCPAddress("0.0.0.0:8888"),
    )

    phello.RegisterHelloServer(s.Server(), &say.Say{})
    phello.RegisterFileServer(s.Server(), &upload.File{})

    ipc.InitIpc(s)

    //imChannel.Channel()
    //imStrem.Stream()
    ImTest.Im()
    s.Run()
}
