package client

import (
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/conn"
    "github.com/go-roc/roc/service/invoke"
    "github.com/gogo/protobuf/proto"
)

func rr(
    c *context.Context,
    s *Client,
    req, rsp proto.Message,
    newInvoke *invoke.Invoke,
) error {
    var cnn *conn.Conn
    var err error

    // if address is nil ,user roundRobin strategy
    // otherwise straight to newInvoke ip server
    if newInvoke.Address() != "" {
        cnn, err = s.strategy.Straight(newInvoke.Scope(), newInvoke.Address())
    } else {
        cnn, err = s.strategy.Next(newInvoke.Scope())
    }
    if err != nil {
        c.Error(err)
        return err
    }

    if newInvoke.FF() {
        newInvoke.InvokeFF(c, req, cnn)
        return nil
    }

    return newInvoke.InvokeRR(c, req, rsp, cnn)
}

func rs(
    c *context.Context,
    s *Client,
    req proto.Message,
    newInvoke *invoke.Invoke,
) chan []byte {
    var cnn *conn.Conn
    var err error

    // if address is nil ,user roundRobin strategy
    // otherwise straight to newInvoke ip server
    if newInvoke.Address() != "" {
        cnn, err = s.strategy.Straight(newInvoke.Scope(), newInvoke.Address())
    } else {
        cnn, err = s.strategy.Next(newInvoke.Scope())
    }

    if err != nil {
        c.Error(err)
        return nil
    }

    //encode req body to roc packet
    b, err := c.Codec().Encode(req)

    if err != nil {
        // create a chan error response
        c.Error(err)
        return nil
    }

    return cnn.Client().RS(c, parcel.Payload(b))
}

func rc(
    c *context.Context,
    s *Client,
    req chan []byte,
    newInvoke *invoke.Invoke,
) chan []byte {
    var cnn *conn.Conn
    var err error

    // if address is nil ,user roundRobin strategy
    // otherwise straight to newInvoke ip server
    if newInvoke.Address() != "" {
        cnn, err = s.strategy.Straight(newInvoke.Scope(), newInvoke.Address())
    } else {
        cnn, err = s.strategy.Next(newInvoke.Scope())
    }
    if err != nil {
        c.Error(err)
        // create a chan error response
        return nil
    }

    return cnn.Client().RC(c, req)
}


