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

package client

import (
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/service/conn"
    "github.com/go-roc/roc/service/invoke"
    "github.com/go-roc/roc/service/opt"
    "github.com/go-roc/roc/service/strategy"
    "github.com/gogo/protobuf/proto"
)

type Client struct {
    //run server option
    opts opt.Option

    //The strategy used by the client to initiate an rpc request, with roundRobin or direct ip request
    strategy strategy.Strategy
}

func NewClient(opts opt.Option) *Client {
    s := &Client{opts: opts}

    s.strategy = strategy.NewStrategy(opts.Endpoint, s.opts.Registry)

    return s
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
        c.Error(err)
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
        c.Error(err)
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
    b, err := c.Codec().Encode(req)

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
