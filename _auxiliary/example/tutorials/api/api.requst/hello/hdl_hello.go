package hello

import (
	"net/http"
	"roc/_auxiliary/example/tutorials/proto/pbhello"
	"roc/client"
	"roc/parcel/context"
	"roc/rlog"
)

type Hello struct {
	opt    client.InvokeOptions
	client pbhello.HelloWorldClient
}

// NewHello new Hello and initialize it for rpc client
// opt is configurable when request.
func NewHello() *Hello {
	return &Hello{client: pbhello.NewHelloWorldClient(client.NewRocClient()), opt: client.WithScope("srv.hello")}
}

func (h *Hello) SayHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	rsp, err := h.client.Say(context.Background(), &pbhello.SayReq{Inc: 1}, h.opt)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	rlog.Info("FROM helloe server: ", rsp.Inc)

	w.Write([]byte("succuess"))
}
