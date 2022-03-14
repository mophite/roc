package rs

import (
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/router"
    "github.com/rsocket/rsocket-go/payload"
    "github.com/rsocket/rsocket-go/rx/flux"
    "github.com/rsocket/rsocket-go/rx/mono"
)

func rr(
    c *context.Context,
    r *router.Router,
    remoteIp string,
    req, rsp *parcel.RocPacket,
) mono.Mono {

    c.RemoteAddr = remoteIp

    err := r.RRProcess(c, req, rsp)

    if err == router.ErrNotFoundHandler {
        c.Errorf("err=%v |path=%s", err, c.Metadata.Method())
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
}

func ff(
    c *context.Context,
    r *router.Router,
    remoteIp string,
    req *parcel.RocPacket,
) {

    c.RemoteAddr = remoteIp

    err := r.FFProcess(c, req)

    if err == router.ErrNotFoundHandler {
        c.Errorf("err=%v |path=%s", err, c.Metadata.Method())
        return
    }

    if err != nil {
        c.Error(err)
        return
    }
}

func rs(
    c *context.Context,
    router *router.Router,
    remoteIp string,
    req *parcel.RocPacket,
    sink flux.Sink,
) {

    c.RemoteAddr = remoteIp

    //if you want to Disconnect channel
    //you must close rsp from server handler
    //this way is very friendly to closing channel transport
    rsp, err := router.RSProcess(c, req)

    //todo cannot know when socket will close to close(rsp)
    //you must close rsp at where send

    if err != nil {
        c.Errorf("transport CC failure |method=%s |err=%v", c.Metadata.Method(), err)
        return
    }

    for b := range rsp {
        data, e := c.Codec().Encode(b)
        if e != nil {
            c.Errorf("transport CC Encode failure |method=%s |err=%v", c.Metadata.Method(), err)
            continue
        }
        sink.Next(payload.New(data, nil))
    }
    sink.Complete()
}

func rc(
    c *context.Context,
    router *router.Router,
    remoteIp string,
    req chan *parcel.RocPacket,
    exitRead chan struct{},
    sink flux.Sink,
) {

    c.RemoteAddr = remoteIp

    //if you want to Disconnect channel
    //you must close rsp from server handler
    //this way is very friendly to closing channel transport
    rsp, err := router.RCProcess(c, req, exitRead)
    if err != nil {
        c.Errorf("transport CC failure |method=%s |err=%v", c.Metadata.Method(), err)
        return
    }

    for b := range rsp {
        data, e := c.Codec().Encode(b)
        if e != nil {
            c.Errorf("transport CC Encode failure |method=%s |err=%v", c.Metadata.Method(), err)
            continue
        }
        sink.Next(payload.New(data, nil))
    }
    sink.Complete()
}
