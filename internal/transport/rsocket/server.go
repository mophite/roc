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
    "sync"

    "github.com/jjeffcaii/reactor-go/scheduler"
    "github.com/rsocket/rsocket-go"
    "github.com/rsocket/rsocket-go/payload"
    "github.com/rsocket/rsocket-go/rx"
    "github.com/rsocket/rsocket-go/rx/flux"
    "github.com/rsocket/rsocket-go/rx/mono"

    "github.com/go-roc/roc/x"

    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/router"

    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/rlog"
)

type server struct {

    //wait server run success
    wg *sync.WaitGroup

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
    r.serverBuilder = rsocket.Receive().OnStart(
        func() {
            r.wg.Done()
        },
    )

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

                var c = context.New()
                var remoteIp, _ = rsocket.GetAddr(sendingSocket)

                if len(r.dog) > 0 {

                    c.SetSetupData(setup.Data())

                    for i := range r.dog {
                        rsp, err := r.dog[i](c)
                        if err != nil {
                            c.Errorf("dog reject you |message=%s", c.Codec().MustEncodeString(rsp))
                            return nil, err
                        }
                    }
                }

                return rsocket.NewAbstractSocket(
                    setupFireAndForget(route, remoteIp, setup),
                    setupRequestResponse(route, remoteIp, setup),
                    setupRequestStream(route, remoteIp, setup),
                    setupRequestChannel(route, remoteIp, r.buffSize, setup),
                ), nil
            },
        )
}

func (r *server) Run(wg *sync.WaitGroup) {
    r.wg = wg
    if r.tcpAddress != "" {
        wg.Add(1)
        r.tcp()
    }

    if r.wssAddress != "" {
        wg.Add(1)
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

func setupRequestResponse(r *router.Router, remoteIp string, setup payload.SetupPayload) rsocket.OptAbstractSocket {
    return rsocket.RequestResponse(
        func(p payload.Payload) mono.Mono {

            c, err := context.FromMetadata(mustGetMetadata(p), setup.DataMimeType(), setup.MetadataMimeType())
            if err != nil {
                rlog.Fatalf("err=%v |metadata=%s |mimeType=%s", err, x.BytesToString(mustGetMetadata(p)), setup.MetadataMimeType())
                return mono.JustOneshot(payload.New(r.Error().Error400(c), nil))
            }

            var req, rsp = parcel.Payload(p.Data()), parcel.NewPacket()

            m := rr(c, r, remoteIp, req, rsp)

            parcel.Recycle(req)
            parcel.Recycle(rsp)

            context.Recycle(c)

            return m
        },
    )
}

func setupFireAndForget(r *router.Router, remoteIp string, setup payload.SetupPayload) rsocket.OptAbstractSocket {
    return rsocket.FireAndForget(
        func(p payload.Payload) {

            var req = parcel.Payload(p.Data())

            c, err := context.FromMetadata(mustGetMetadata(p), setup.DataMimeType(), setup.MetadataMimeType())
            if err != nil {
                rlog.Fatalf("err=%v |metadata=%s |mimeType=%s", err, x.BytesToString(mustGetMetadata(p)), setup.MetadataMimeType())
                return
            }

            ff(c, r, remoteIp, req)

            parcel.Recycle(req)

            context.Recycle(c)
        },
    )
}

func (r *server) Close() {
    return
}

func setupRequestStream(router *router.Router, remoteIp string, setup payload.SetupPayload) rsocket.OptAbstractSocket {
    return rsocket.RequestStream(
        func(p payload.Payload) flux.Flux {

            return flux.Create(
                func(ctx ctx.Context, sink flux.Sink) {

                    var req = parcel.Payload(p.Data())

                    c, err := context.FromMetadata(mustGetMetadata(p), setup.DataMimeType(), setup.MetadataMimeType())
                    if err != nil {
                        rlog.Fatalf("err=%v |metadata=%s |mimeType=%s", err, x.BytesToString(mustGetMetadata(p)), setup.MetadataMimeType())
                        return
                    }

                    rs(c, router, remoteIp, req, sink)

                    parcel.Recycle(req)
                },
            )
        },
    )
}

func setupRequestChannel(router *router.Router, remoteIp string, buffSize int, setup payload.SetupPayload) rsocket.OptAbstractSocket {
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
                        rlog.Debugf("requestChanel success |ip=%s |meta=%s", remoteIp, string(meta))
                        break
                    }

                    c, err := context.FromMetadata(meta, setup.DataMimeType(), setup.MetadataMimeType())
                    if err != nil {
                        rlog.Errorf("err=%v |metadata=%s |mimeType=%s", err, x.BytesToString(meta), setup.MetadataMimeType())
                        return
                    }

                    rc(c, router, remoteIp, req, exitRead, sink)
                },
            )
        },
    )
}
