package rs

import (
	ctx "context"
	"github.com/jjeffcaii/reactor-go/scheduler"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"time"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"

	"roc/internal/endpoint"
	"roc/parcel"
	"roc/parcel/context"
	"roc/rlog"
)

type client struct {
	client rsocket.Client

	connectTimeout    time.Duration
	keepaliveInterval time.Duration
	keepaliveLifetime time.Duration
}

func NewClient(connTimeout, interval, tll time.Duration) *client {
	return &client{
		connectTimeout:    connTimeout,
		keepaliveInterval: interval,
		keepaliveLifetime: tll,
	}
}

func (cli *client) Dial(e *endpoint.Endpoint, ch chan string) (err error) {
	cli.client, err = rsocket.
		Connect().
		Scheduler(scheduler.Elastic(), nil).
		KeepAlive(cli.keepaliveInterval, cli.keepaliveLifetime, 1).
		ConnectTimeout(cli.connectTimeout).
		OnConnect(func(client rsocket.Client, err error) {
			rlog.Debugf("connected at: %s", e.Address)
		}).
		OnClose(func(err error) {
			rlog.Debugf("server [%s %s] is closed |err=%v", e.Name, e.Address, err)
			ch <- e.Address
		}).
		Transport(rsocket.TCPClient().SetAddr(e.Address).Build()).
		Start(ctx.TODO())
	return err
}

func (cli *client) RR(c *context.Context, req *parcel.RocPacket, rsp *parcel.RocPacket) (err error) {
	pl, release, err := cli.
		client.
		RequestResponse(payload.New(req.Bytes(), c.Body())).
		BlockUnsafe(ctx.Background())

	if err != nil {
		c.Error("socket err occurred ", err)
		return err
	}

	rsp.Write(pl.Data())

	release()

	return nil
}

func (cli *client) RS(c *context.Context, req *parcel.RocPacket) (chan []byte, chan error) {
	var (
		f    = cli.client.RequestStream(payload.New(req.Bytes(), c.Body()))
		rsp  = make(chan []byte)
		errs = make(chan error)
	)

	f.
		SubscribeOn(scheduler.Parallel()).
		DoFinally(func(s rx.SignalType) {
			//todo handler rx.SignalType
			close(rsp)
			close(errs)
		}).
		Subscribe(
			ctx.Background(),
			rx.OnNext(func(p payload.Payload) error {
				rsp <- payload.Clone(p).Data()
				return nil
			}),
			rx.OnError(func(e error) {
				errs <- e
			}),
		)

	parcel.Recycle(req)

	return rsp, errs
}

func (cli *client) RC(c *context.Context, req chan []byte, errIn chan error) (chan []byte, chan error) {
	var (
		sendPayload = make(chan payload.Payload, cap(req))
	)

	go func() {
		sendPayload <- payload.New(c.Body(), nil)
		for d := range req {
			pl := payload.New(d, nil)
			sendPayload <- pl
		}

		close(sendPayload)
	}()

	var (
		f    = cli.client.RequestChannel(flux.CreateFromChannel(sendPayload, errIn))
		rsp  = make(chan []byte)
		errs = make(chan error)
	)

	f.
		SubscribeOn(scheduler.Parallel()).
		DoFinally(func(s rx.SignalType) {
			//todo handler rx.SignalType
			close(rsp)
			close(errs)
		}).
		Subscribe(
			ctx.Background(),
			rx.OnNext(func(p payload.Payload) error {
				rsp <- payload.Clone(p).Data()
				return nil
			}),
			rx.OnError(func(e error) {
				errs <- e
			}),
		)

	return rsp, errs
}

func (cli *client) String() string {
	return "rsocket"
}

func (cli *client) Close() {
	_ = cli.client.Close()
}
