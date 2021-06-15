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
	"github.com/go-roc/roc/parcel/codec"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/internal/transport"
	"github.com/go-roc/roc/parcel/context"
)

type RocClient struct {
	opts option

	//The strategy used by the client to initiate an rpc request, with roundRobin or direct ip request
	strategy Strategy
}

// Registry Get the use of serviceName discovery
func (r *RocClient) Registry() registry.Registry {
	return r.opts.registry
}

func (r *RocClient) Client() transport.Client {
	return r.opts.client
}

func (r *RocClient) ConnectTimeout() time.Duration {
	return r.opts.connectTimeout
}

func (r *RocClient) KeepaliveInterval() time.Duration {
	return r.opts.keepaliveInterval
}

func (r *RocClient) KeepaliveLifetime() time.Duration {
	return r.opts.keepaliveLifetime
}

func (r *RocClient) Codec() codec.Codec {
	return r.opts.cc
}

// NewRocClient create a client
// set opts,or use default
func NewRocClient(opts ...Options) (lp *RocClient) {

	lp = &RocClient{opts: newOpts(opts...)}
	lp.strategy = newStrategy(lp.opts.registry, lp.opts.client)

	return
}

// InvokeRR rpc request requestResponse,it's block request,one request one response
//
func (r *RocClient) InvokeRR(c *context.Context, method string, req, rsp proto.Message, opts ...InvokeOptions) error {

	// new a invoke setting
	invoke, err := newInvoke(c, method, r, opts...)
	if err != nil {
		return err
	}

	var conn *Conn

	// if address is nil ,user roundRobin strategy
	// otherwise straight to invoke ip server
	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}
	if err != nil {
		return err
	}

	return invoke.invokeRR(c, req, rsp, conn, r.opts)
}

// InvokeRS rpc request requestStream,it's one request and multiple response
//
func (r *RocClient) InvokeRS(c *context.Context, method string, req proto.Message, opts ...InvokeOptions) (chan []byte, chan error) {

	// new a invoke setting
	invoke, err := newInvoke(c, method, r, opts...)
	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	var conn *Conn

	// if address is nil ,user roundRobin strategy
	// otherwise straight to invoke ip server
	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}

	//encode req body to roc packet
	b, err := r.opts.cc.Encode(req)

	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	return conn.Client().RS(c, parcel.Payload(b))
}

func (r *RocClient) InvokeRC(c *context.Context, method string, req chan []byte, errIn chan error, opts ...InvokeOptions) (chan []byte, chan error) {

	// new a invoke setting
	invoke, err := newInvoke(c, method, r, opts...)
	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	var conn *Conn

	// if address is nil ,user roundRobin strategy
	// otherwise straight to invoke ip server
	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}
	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	return conn.Client().RC(c, req, errIn)
}

// Close handler client close
func (r *RocClient) Close() {
	if r.strategy != nil {
		r.strategy.Close()
	}
}
