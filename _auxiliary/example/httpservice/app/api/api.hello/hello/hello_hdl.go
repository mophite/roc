package hello

import (
	"github.com/go-roc/roc/_auxiliary/example/httpservice/proto/phello"
	"github.com/go-roc/roc/parcel/context"
)

type Hello struct {
}

func (h *Hello) Say(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) {
	c.Info("--------hello--------", req.Ping)
	rsp.Pong = "pong"
}
