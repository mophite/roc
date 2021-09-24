package ipc

import (
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

var sayClient phello.HelloSrvClient

func SaySrv(c *context.Context, req *phello.SayReq) (rsp *phello.SayRsp, err error) {
    return sayClient.SaySrv(c, req, invokeHello)
}

func SayChannel(c *context.Context, req chan *phello.SayReq) chan *phello.SayRsp {
    return sayClient.SayChannel(c, req, invokeHello)
}

func SayStream(c *context.Context, req *phello.SayReq) chan *phello.SayRsp {
    return sayClient.SayStream(c, req, invokeHello)
}
