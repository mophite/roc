package client

import (
    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/codec"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/conn"
    "github.com/go-roc/roc/service/invoke"
    "github.com/go-roc/roc/service/strategy"
    "github.com/gogo/protobuf/proto"
)

type Client struct {
    //run server option
    opts Option

    //The strategy used by the client to initiate an rpc request, with roundRobin or direct ip request
    strategy strategy.Strategy
}

func NewClient(opts ...Options) *Client {
    s := &Client{opts: newOpts(opts...)}

    s.strategy = strategy.NewStrategy(endpoint.DefaultLocalEndpoint, s.opts.registry)

    return s
}

func (s *Client) Codec() codec.Codec {
    return s.opts.cc
}

// InvokeRR rpc request requestResponse,it's block request,one request one response
func (s *Client) InvokeRR(
    c *context.Context,
    method string,
    req, rsp proto.Message,
    opts ...invoke.InvokeOptions,
) error {

    // new a newInvoke setting
    newInvoke, err := invoke.NewInvoke(c, method, opts...)
    if err != nil {
        return err
    }

    var cnn *conn.Conn

    // if address is nil ,user roundRobin strategy
    // otherwise straight to newInvoke ip server
    if newInvoke.Address() != "" {
        cnn, err = s.strategy.Straight(newInvoke.Scope(), newInvoke.Address())
    } else {
        cnn, err = s.strategy.Next(newInvoke.Scope())
    }

    if err != nil {
        return err
    }

    return newInvoke.InvokeRR(c, req, rsp, cnn)
}

// InvokeRS rpc request requestStream,it's one request and multiple response
func (s *Client) InvokeRS(c *context.Context, method string, req proto.Message, opts ...invoke.InvokeOptions) (
    chan []byte,
    chan error,
) {

    // new a newInvoke setting
    newInvoke, err := invoke.NewInvoke(c, method)
    if err != nil {
        // create a chan error response
        var errs = make(chan error)
        errs <- err
        close(errs)
        return nil, errs
    }

    var cnn *conn.Conn

    // if address is nil ,user roundRobin strategy
    // otherwise straight to newInvoke ip server
    if newInvoke.Address() != "" {
        cnn, err = s.strategy.Straight(newInvoke.Scope(), newInvoke.Address())
    } else {
        cnn, err = s.strategy.Next(newInvoke.Scope())
    }

    //encode req body to roc packet
    b, err := s.opts.cc.Encode(req)

    if err != nil {
        // create a chan error response
        var errs = make(chan error)
        errs <- err
        close(errs)
        return nil, errs
    }

    return cnn.Client().RS(c, parcel.Payload(b))
}

// InvokeRC rpc request requestChannel,it's multiple request and multiple response
func (s *Client) InvokeRC(
    c *context.Context,
    method string,
    req chan []byte,
    errIn chan error,
    opts ...invoke.InvokeOptions,
) (chan []byte, chan error) {

    // new a newInvoke setting
    newInvoke, err := invoke.NewInvoke(c, method, opts...)
    if err != nil {
        // create a chan error response
        var errs = make(chan error)
        errs <- err
        close(errs)
        return nil, errs
    }

    var cnn *conn.Conn

    // if address is nil ,user roundRobin strategy
    // otherwise straight to newInvoke ip server
    if newInvoke.Address() != "" {
        cnn, err = s.strategy.Straight(newInvoke.Scope(), newInvoke.Address())
    } else {
        cnn, err = s.strategy.Next(newInvoke.Scope())
    }
    if err != nil {
        // create a chan error response
        var errs = make(chan error)
        errs <- err
        close(errs)
        return nil, errs
    }

    return cnn.Client().RC(c, req, errIn)
}

func (s *Client) Close() {
    if s.strategy != nil {
        s.strategy.Close()
    }
}
