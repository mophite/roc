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
    "fmt"
    "time"

    "github.com/jjeffcaii/reactor-go/scheduler"
    "github.com/rsocket/rsocket-go"
    "github.com/rsocket/rsocket-go/rx"
    "github.com/rsocket/rsocket-go/rx/flux"

    "github.com/rsocket/rsocket-go/payload"

    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/rlog"
)

//this is rsocket client
type client struct {

    //rsocket client
    client rsocket.Client

    //rsocket connect timeout
    connectTimeout time.Duration

    //rsocket keepalive interval
    keepaliveInterval time.Duration

    //rsocket keepalive life time
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
        Scheduler(scheduler.Elastic(), nil). //set scheduler to best
        KeepAlive(cli.keepaliveInterval, cli.keepaliveLifetime, 1).
        ConnectTimeout(cli.connectTimeout).
        OnConnect(
            func(client rsocket.Client, err error) { //handler when connect success
                rlog.Debugf("connected at: %s", e.Address)
            },
        ).
        OnClose(
            func(err error) { //when net occur some error,it's will be callback the error server ip address
                rlog.Debugf("server [%s %s] is closed |err=%v", e.Name, e.Address, err)
                ch <- e.Address
            },
        ).
        Transport(rsocket.TCPClient().SetAddr(e.Address).Build()). //setup transport and start
        Start(ctx.TODO())
    return err
}

// RR requestResponse on blockUnsafe
func (cli *client) RR(c *context.Context, req *parcel.RocPacket, rsp *parcel.RocPacket) (err error) {
    fmt.Println("-----now", req.Len())
    now := time.Now()
    pl, err := cli.
        client.
        RequestResponse(payload.New(req.Bytes(), c.Body())).
        Block(ctx.Background())

    if err != nil {
        c.Error("socket err occurred ", err)
        return err
    }
    fmt.Println("------now latency----", time.Since(now).Seconds())

    rsp.Write(pl.Data())

    //release()

    return nil
}

// RS requestStream
func (cli *client) RS(c *context.Context, req *parcel.RocPacket) (chan []byte, chan error) {
    var (
        f    = cli.client.RequestStream(payload.New(req.Bytes(), c.Body()))
        rsp  = make(chan []byte)
        errs = make(chan error)
    )

    f.
        SubscribeOn(scheduler.Parallel()).
        DoFinally(
            func(s rx.SignalType) {
                //todo handler rx.SignalType
                close(rsp)
                close(errs)
            },
        ).
        Subscribe(
            ctx.Background(),
            rx.OnNext(
                func(p payload.Payload) error {
                    rsp <- payload.Clone(p).Data()
                    return nil
                },
            ),
            rx.OnError(
                func(e error) {
                    errs <- e
                },
            ),
        )

    parcel.Recycle(req)

    return rsp, errs
}

// RC requestChannel
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
        DoFinally(
            func(s rx.SignalType) {
                //todo handler rx.SignalType
                close(rsp)
                close(errs)
            },
        ).
        Subscribe(
            ctx.Background(),
            rx.OnNext(
                func(p payload.Payload) error {
                    rsp <- payload.Clone(p).Data()
                    return nil
                },
            ),
            rx.OnError(
                func(e error) {
                    errs <- e
                },
            ),
        )

    return rsp, errs
}

func (cli *client) String() string {
    return "rsocket"
}

func (cli *client) Close() {
    if cli.client != nil {
        _ = cli.client.Close()
    }
}
