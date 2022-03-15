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
    "github.com/go-roc/roc/parcel/context"
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
    cc, newInvoke, err := invoke.NewInvoke(c, method, opts...)
    if err != nil {
        c.Error(err)
        return err
    }

    err = rr(cc, s, req, rsp, newInvoke)

    //context.Recycle(cc)

    return err
}

// InvokeRS rpc request requestStream,it's one request and multiple response
func (s *Client) InvokeRS(
    c *context.Context,
    method string,
    req proto.Message,
    opts ...invoke.InvokeOptions,
) chan []byte {

    // new a newInvoke setting
    cc, newInvoke, err := invoke.NewInvoke(c, method, opts...)
    if err != nil {
        // create a chan error response
        c.Error(err)
        return nil
    }

    return rs(cc, s, req, newInvoke)
}

// InvokeRC rpc request requestChannel,it's multiple request and multiple response
func (s *Client) InvokeRC(
    c *context.Context,
    method string,
    req chan []byte,
    opts ...invoke.InvokeOptions,
) chan []byte {

    // new a newInvoke setting
    cc, newInvoke, err := invoke.NewInvoke(c, method, opts...)
    if err != nil {
        c.Error(err)
        // create a chan error response
        return nil
    }

    return rc(cc, s, req, newInvoke)
}

func (s *Client) CloseClient() {
    if s.strategy != nil {
        s.strategy.CloseStrategy()
    }
}
