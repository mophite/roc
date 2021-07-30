package hello

import (
    "fmt"

    "github.com/go-roc/roc/_auxiliary/example/httpservice/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

type Hello struct {
}

func (h *Hello) Say(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) (err error) {
    c.Info("--------hello--------", req.Ping)
    rsp.Pong = "pong"

    return nil
}

func (h *Hello) SayGet(c *context.Context, req *phello.ApiReq, rsp *phello.ApiRsp) (err error) {
    fmt.Println("------get------", req.Params["name"])
    rsp.Code = 200
    rsp.Msg = "success"
    return err
}
