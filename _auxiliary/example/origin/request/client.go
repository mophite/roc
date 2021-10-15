package main

import (
    ctx "context"
    "runtime"
    "time"

    "github.com/go-roc/roc/rlog"
    "github.com/jjeffcaii/reactor-go/scheduler"
    "github.com/rsocket/rsocket-go"
)

func newClient() rsocket.Client {
    client, err := rsocket.
        Connect().
        Scheduler(
            scheduler.NewElastic(runtime.NumCPU()<<8),
            scheduler.NewElastic(runtime.NumCPU()<<8),
        ). //set scheduler to best
        KeepAlive(time.Second*5, time.Second*5, 1).
        ConnectTimeout(time.Second * 5).
        OnConnect(
            func(client rsocket.Client, err error) { //handler when connect success
                rlog.Debug("connected success")
            },
        ).
        OnClose(
            func(err error) { //when net occur some error,it's will be callback the error server ip address
                rlog.Error(err)
            },
        ).
        Transport(rsocket.TCPClient().SetAddr("0.0.0.0:8888").Build()). //setup transport and start
        Start(ctx.TODO())
    if err != nil {
        panic(err)
    }
    return client
}
