package ipc

import (
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

var imClient phello.ImClient

func Connect(c *context.Context, req *phello.ConnectReq) (rsp *phello.ConnectRsp, err error) {
    return imClient.Connect(c, req, invokeHello)
}

func Count(c *context.Context, req *phello.CountReq) (rsp *phello.CountRsp, err error) {
    return imClient.Count(c, req, invokeHello)
}

func SendMessage(c *context.Context, req chan *phello.SendMessageReq) (
    rsp chan *phello.SendMessageRsp, exit chan struct{},
) {
    return imClient.SendMessage(c, req, invokeHello)
}
