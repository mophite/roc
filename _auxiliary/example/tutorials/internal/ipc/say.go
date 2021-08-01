package ipc

import (
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/invoke"
)

var sayClient phello.HelloSrvClient

func SaySrv(c *context.Context, req *phello.SayReq) (rsp *phello.SayRsp, err error) {
    return sayClient.SaySrv(c, req, invoke.WithName("srv.hello"))
}

func SayChannel(c *context.Context, req chan *phello.SayReq, errsIn chan error) (
    rsp chan *phello.SayRsp,
    err chan error,
) {
    return sayClient.SayChannel(c, req, errsIn, invoke.WithName("srv.hello"))
}

func SayStream(c *context.Context, req *phello.SayReq) (
    rsp chan *phello.SayRsp,
    err chan error,
) {
    return sayClient.SayStream(c, req, invoke.WithName("srv.hello"))
}
