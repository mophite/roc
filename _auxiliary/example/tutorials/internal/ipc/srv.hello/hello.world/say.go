package hello_world

import (
	"tutorials/proto/phello"

	"github.com/go-roc/roc"
)

var (
	//rpc srv.hello HelloHelloWorld invoke options
	srvHelloHelloWorldOpt = []roc.InvokeOptions{roc.WithName("srv.hello")}
)

//SayHello is rpc service what ...
func SayHello(c *context.Context, client phello.HelloWorldClient, req *phello.SayReq) (rsp *phello.SayRsp, err error) {
	//can handle other business
	return client.Say(c, req, srvHelloHelloWorldOpt...)
}
