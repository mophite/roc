package rs

import (
	rc "context"
	"github.com/jjeffcaii/reactor-go/scheduler"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"roc/internal/router"
	"roc/parcel"
	"roc/parcel/context"
	"roc/rlog"
	"runtime"
)

type server struct {
	serverName    string
	tcpAddress    string
	wssAddress    string
	buffSize      int
	serverBuilder rsocket.ServerBuilder
	serverStart   rsocket.ToServerStarter
}

func (r *server) Address() string {
	return "[tcp: " + r.tcpAddress + "] [wss: " + r.wssAddress + "]"
}

func (r *server) String() string {
	return "rsocket"
}

func NewServer(tcpAddress, wssAddress, serverName string, buffSize int) *server {
	return &server{
		serverName: serverName,
		tcpAddress: tcpAddress,
		wssAddress: wssAddress,
		buffSize:   buffSize,
	}
}

func (r *server) Accept(route *router.Router) {
	r.serverBuilder = rsocket.Receive()

	r.serverBuilder.
		Scheduler(nil, scheduler.NewElastic(runtime.NumCPU()*2))

	r.serverBuilder.
		OnStart(func() {
			// todo
		})

	r.serverBuilder.Resume()

	r.serverStart = r.serverBuilder.
		Acceptor(func(
			ctx rc.Context,
			setup payload.SetupPayload,
			sendingSocket rsocket.CloseableRSocket,
		) (rsocket.RSocket, error) {
			return rsocket.NewAbstractSocket(
				setupRequestResponse(route),
				setupRequestStream(route),
				setupRequestChannel(route, r.buffSize),
			), nil
		})
}

func (r *server) Run() {
	if r.tcpAddress != "" {
		r.tcp()
	}

	if r.wssAddress != "" {
		r.wss()
	}
}

func (r *server) tcp() {
	go func() {
		err := r.serverStart.Transport(
			rsocket.
				TCPServer().
				SetAddr(r.tcpAddress).
				Build()).Serve(rc.TODO())

		if err != nil {
			rlog.Errorf("tcp server start err=%v", err)
		}
	}()

}

func (r *server) wss() {
	go func() {
		err := r.serverStart.Transport(
			rsocket.
				WebsocketServer().
				SetAddr(r.wssAddress).
				Build()).Serve(rc.TODO())

		if err != nil {
			rlog.Errorf("wss server start err=%v", err)
		}
	}()
}

func getMetadata(p payload.Payload) []byte {
	b, _ := p.Metadata()
	return b
}

func setupRequestResponse(router *router.Router) rsocket.OptAbstractSocket {
	return rsocket.RequestResponse(func(p payload.Payload) mono.Mono {

		var req, rsp = parcel.Payload(p.Data()), parcel.NewPacket()
		defer func() {
			parcel.Recycle(req, rsp)
		}()

		err := router.RRProcess(context.FromMetadata(getMetadata(p)), req, rsp)
		if err != nil {
			return mono.JustOneshot(
				payload.New(router.Error().
					Encode(parcel.ErrorCodeBadRequest, err), nil))
		}

		return mono.JustOneshot(payload.New(rsp.Bytes(), nil))
	})
}

func (r *server) Close() {
	return
}

func setupRequestStream(router *router.Router) rsocket.OptAbstractSocket {
	return rsocket.RequestStream(func(p payload.Payload) flux.Flux {

		var req = parcel.Payload(p.Data())

		rsp, errs := router.RSProcess(context.FromMetadata(getMetadata(p)), req)

		parcel.Recycle(req)

		f := flux.Create(func(ctx rc.Context, sink flux.Sink) {
		QUIT:
			for {
				select {
				case b, ok := <-rsp:
					if ok {
						data, err := router.Codec().Encode(b)
						if err != nil {
							rlog.Error(err)
							break
						}
						sink.Next(payload.New(data, nil))
					} else {
						break QUIT
					}
				case e := <-errs:
					if e != nil {
						rlog.Error(e)
						break QUIT
					}
				}
			}

			sink.Complete()
		})

		return f
	})
}

func setupRequestChannel(router *router.Router, buffSize int) rsocket.OptAbstractSocket {
	return rsocket.RequestChannel(func(f flux.Flux) flux.Flux {
		var (
			errs = make(chan error)
			req  = make(chan *parcel.RocPacket, buffSize)
		)

		f.SubscribeOn(scheduler.Parallel()).
			DoFinally(func(s rx.SignalType) {
				//todo handler rx.SignalType
				close(req)
				close(errs)
			}).
			Subscribe(
				rc.Background(),
				rx.OnNext(func(p payload.Payload) error {
					req <- parcel.Payload(payload.Clone(p).Data())
					return nil
				}),
				rx.OnError(func(e error) {
					errs <- e
				}),
			)

		return flux.Create(func(ctx rc.Context, sink flux.Sink) {

			var meta []byte
			for b := range req {
				meta = b.Bytes()
				break
			}

			rsp, outErrs := router.RCProcess(context.FromMetadata(meta), req, errs)

		QUIT:
			for {
				select {
				case b, ok := <-rsp:
					if ok {
						data, e := router.Codec().Encode(b)
						if e != nil {
							rlog.Error(e)
							break
						}
						sink.Next(payload.New(data, nil))
					} else {
						break QUIT
					}
				case e := <-outErrs:
					if e != nil {
						rlog.Error(e)
						break QUIT
					}
				}
			}
			sink.Complete()
		})
	})
}
