package hello

import (
	"github.com/go-roc/roc/_auxiliary/example/hello/proto/phello"
	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/rlog"
	"github.com/go-roc/roc/service/client"
	"github.com/go-roc/roc/service/invoke"
)

type Hello struct {
	Client *client.Client
}

var sayClient phello.HelloWorldClient

func (h *Hello) Say(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) error {
	c.Info("--------api hello--------", req.Ping, c.ContentType)
	if sayClient == nil {
		sayClient = phello.NewHelloWorldClient(h.Client)
	}

	sayRsp, err := sayClient.Say(c, req, invoke.WithName("srv.hello"))
	if err != nil {
		rlog.Error(err)
		rsp.Pong = "error"
		return err
	}

	rsp.Pong = sayRsp.Pong

	return nil
}
