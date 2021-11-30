// Copyright (c) 2021 roc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package rs

import (
    ctx "context"
    "runtime"

    "github.com/jjeffcaii/reactor-go/scheduler"
    "github.com/rsocket/rsocket-go"
    "github.com/rsocket/rsocket-go/payload"
    "github.com/rsocket/rsocket-go/rx"
    "github.com/rsocket/rsocket-go/rx/flux"
    "github.com/rsocket/rsocket-go/rx/mono"

    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/router"

    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/rlog"
)

type server struct {

    //given serverName to service discovery to find
    serverName string

    //tcp socket address
    tcpAddress string

    //websocket address
    wssAddress string

    //requestChannel buffSize setting
    buffSize int

    //rsocket serverBuilder
    serverBuilder rsocket.ServerBuilder

    //rsocket serverStarter
    serverStart rsocket.ToServerStarter

    dog []handler.DogHandler
}

func (r *server) Address() string {
    return "[tcp: " + r.tcpAddress + "] [wss: " + r.wssAddress + "]"
}

func (r *server) String() string {
    return "rsocket"
}

func NewServer(tcpAddress, wssAddress, serverName string, buffSize int, dog ...handler.DogHandler) *server {
    return &server{
        serverName: serverName,
        tcpAddress: tcpAddress,
        wssAddress: wssAddress,
        buffSize:   buffSize,
        dog:        dog,
    }
}

func (r *server) Accept(route *router.Router) {
    r.serverBuilder = rsocket.Receive()

    r.serverBuilder.Scheduler(
        scheduler.NewElastic(runtime.NumCPU()<<8),
        scheduler.NewElastic(runtime.NumCPU()<<8),
    ) // setting scheduler goroutine on numCPU*2 to better working

    r.serverBuilder.Resume()
    r.serverStart = r.serverBuilder.
        Acceptor(
            func(
                cc ctx.Context,
                setup payload.SetupPayload,
                sendingSocket rsocket.CloseableRSocket,
            ) (rsocket.RSocket, error) {

                var c = context.Background()

                if len(r.dog) > 0 {

                    c.SetSetupData(setup.Data())
                    c.RemoteAddr, _ = rsocket.GetAddr(sendingSocket)

                    for i := range r.dog {
                        rsp, err := r.dog[i](c)
                        if err != nil {
                            c.Errorf("dog reject you |message=%s", c.Codec().MustEncodeString(rsp))
                            return nil, err
                        }
                    }
                }

                return rsocket.NewAbstractSocket(
                    setupRequestResponse(route),
                    setupRequestStream(route),
                    setupRequestChannel(route, r.buffSize),
                ), nil
            },
        )
}

func (r *server) Run() {
    if r.tcpAddress != "" {
        r.tcp()
    }

    if r.wssAddress != "" {
        r.wss()
    }
}

//run tcp socket server
func (r *server) tcp() {
    go func() {
        err := r.serverStart.Transport(
            rsocket.
                TCPServer().
                SetAddr(r.tcpAddress).
                Build(),
        ).Serve(ctx.TODO())

        if err != nil {
            panic(err)
        }
    }()
}

//run websocket server
func (r *server) wss() {
    go func() {
        err := r.serverStart.Transport(
            rsocket.
                WebsocketServer().
                SetAddr(r.wssAddress).
                Build(),
        ).Serve(ctx.TODO())

        if err != nil {
            panic(err)
        }
    }()
}

// get metadata ignore error
func mustGetMetadata(p payload.Payload) []byte {
    b, _ := p.Metadata()
    return b
}

func setupRequestResponse(r *router.Router) rsocket.OptAbstractSocket {
    return rsocket.RequestResponse(
        func(p payload.Payload) mono.Mono {

            var req, rsp = parcel.Payload(p.Data()), parcel.NewPacket()
            defer func() {
                parcel.Recycle(req, rsp)
            }()

            var c = context.FromMetadata(mustGetMetadata(p))

            err := r.RRProcess(c, req, rsp)

            if err == router.ErrNotFoundHandler {
                c.Errorf("err=%v |path=%s", err, c.Method())
                return mono.JustOneshot(payload.New(r.Error().Error404(c), nil))
            }

            if err != nil && rsp.Len() > 0 {
                c.Error(err)
                return mono.JustOneshot(payload.New(rsp.Bytes(), nil))
            }

            if err != nil {
                c.Error(err)
                return mono.JustOneshot(payload.New(r.Error().Error400(c), nil))
            }

            return mono.JustOneshot(payload.New(rsp.Bytes(), nil))
        },
    )
}

func (r *server) Close() {
    return
}

func setupRequestStream(router *router.Router) rsocket.OptAbstractSocket {
    return rsocket.RequestStream(
        func(p payload.Payload) flux.Flux {

            return flux.Create(
                func(ctx ctx.Context, sink flux.Sink) {
                    var (
                        req = parcel.Payload(p.Data())
                    )

                    var c = context.FromMetadata(mustGetMetadata(p))

                    //if you want to Disconnect channel
                    //you must close rsp from server handler
                    //this way is very friendly to closing channel transport
                    rsp, err := router.RSProcess(c, req)

                    //todo cannot know when socket will close to close(rsp)
                    //you must close rsp at where send

                    parcel.Recycle(req)

                    if err != nil {
                        c.Errorf("transport CC failure |method=%s |err=%v", c.Method(), err)
                        return
                    }

                    for b := range rsp {
                        data, e := c.Codec().Encode(b)
                        if e != nil {
                            c.Errorf("transport CC Encode failure |method=%s |err=%v", c.Method(), err)
                            continue
                        }
                        sink.Next(payload.New(data, nil))
                    }
                    sink.Complete()
                },
            )
        },
    )
}

func setupRequestChannel(router *router.Router, buffSize int) rsocket.OptAbstractSocket {
    return rsocket.RequestChannel(
        func(f flux.Flux) flux.Flux {
            var (
                req      = make(chan *parcel.RocPacket, buffSize)
                exitRead = make(chan struct{})
            )

            //read data from client by channel transport method
            f.SubscribeOn(scheduler.Parallel()).
                DoFinally(
                    func(s rx.SignalType) {
                        close(req)
                        close(exitRead)
                    },
                ).
                Subscribe(
                    ctx.Background(),
                    rx.OnNext(
                        func(p payload.Payload) error {
                            req <- parcel.Payload(payload.Clone(p).Data())
                            return nil
                        },
                    ),
                    rx.OnError(
                        //if client is occurred error
                        func(e error) {
                            rlog.Errorf("setupRequestChannel OnError |err=%v", e)
                        },
                    ),
                )

            return flux.Create(
                func(ctx ctx.Context, sink flux.Sink) {

                    var meta []byte
                    for b := range req {
                        meta = b.Bytes()
                        break
                    }

                    var c = context.FromMetadata(meta)

                    //if you want to Disconnect channel
                    //you must close rsp from server handler
                    //this way is very friendly to closing channel transport
                    rsp, err := router.RCProcess(c, req, exitRead)
                    if err != nil {
                        c.Errorf("transport CC failure |method=%s |err=%v", c.Method(), err)
                        return
                    }

                    for b := range rsp {
                        data, e := c.Codec().Encode(b)
                        if e != nil {
                            c.Errorf("transport CC Encode failure |method=%s |err=%v", c.Method(), err)
                            continue
                        }
                        sink.Next(payload.New(data, nil))
                    }
                    sink.Complete()
                },
            )
        },
    )
}
