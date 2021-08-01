package ipc

import (
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/invoke"
)

var imClient phello.ImClient

func Connect(c *context.Context, req *phello.ConnectReq) (rsp *phello.ConnectRsp, err error) {
    return imClient.Connect(c, req, invoke.WithName("srv.hello"))
}

func Count(c *context.Context, req *phello.CountReq) (rsp *phello.CountRsp, err error) {
    return imClient.Count(c, req, invoke.WithName("srv.hello"))
}

func SendMessage(c *context.Context, req chan *phello.SendMessageReq, errsIn chan error) (
    rsp chan *phello.SendMessageRsp,
    err chan error,
) {
    return imClient.SendMessage(c, req, errsIn, invoke.WithName("srv.hello"))
}
