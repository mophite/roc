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

    "github.com/go-roc/roc/internal/router"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
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

    r.serverBuilder.Scheduler(
        nil,
        scheduler.NewElastic(runtime.NumCPU()*2),
    ) // setting scheduler goroutine on numCPU*2 to better working
    //
    r.serverBuilder.Resume()

    r.serverStart = r.serverBuilder.
        Acceptor(
            func(
                ctx ctx.Context,
                setup payload.SetupPayload,
                sendingSocket rsocket.CloseableRSocket,
            ) (rsocket.RSocket, error) {
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

func setupRequestResponse(router *router.Router) rsocket.OptAbstractSocket {
    return rsocket.RequestResponse(
        func(p payload.Payload) mono.Mono {

            var req, rsp = parcel.Payload(p.Data()), parcel.NewPacket()
            defer func() {
                parcel.Recycle(req, rsp)
            }()

            err := router.RRProcess(context.FromMetadata(mustGetMetadata(p)), req, rsp)
            if err != nil {
                return mono.JustOneshot(
                    payload.New(
                        router.Error().
                            Encode(parcel.ErrorCodeBadRequest, err), nil,
                    ),
                )
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

            var req = parcel.Payload(p.Data())

            rsp, errs := router.RSProcess(context.FromMetadata(mustGetMetadata(p)), req)

            parcel.Recycle(req)

            f := flux.Create(
                func(ctx ctx.Context, sink flux.Sink) {
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
                },
            )

            return f
        },
    )
}

func setupRequestChannel(router *router.Router, buffSize int) rsocket.OptAbstractSocket {
    return rsocket.RequestChannel(
        func(f flux.Flux) flux.Flux {
            var (
                errs = make(chan error)
                req  = make(chan *parcel.RocPacket, buffSize)
            )

            f.SubscribeOn(scheduler.Parallel()).
                DoFinally(
                    func(s rx.SignalType) {
                        //todo handler rx.SignalType
                        close(req)
                        close(errs)
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
                        func(e error) {
                            errs <- e
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
                },
            )
        },
    )
}
