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

package invoke

import (
	"errors"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/service/conn"
	"github.com/go-roc/roc/x/backoff"

	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/parcel"
)

type Invoke struct {
	// invoke options
	opts InvokeOption
}

// NewInvoke create a invoke
func NewInvoke(c *context.Context, method string, opts ...InvokeOptions) (*context.Context, *Invoke, error) {
	invoke := &Invoke{}

	for i := range opts {
		opts[i](&invoke.opts)
	}

	if invoke.opts.serviceName == "" || invoke.opts.scope == "" {
		return nil, nil, errors.New("not set rpc service name")
	}

	method = invoke.opts.prefix + method

	// initialize tunnel for requestChannel only
	if invoke.opts.buffSize == 0 {
		invoke.opts.buffSize = 10
	}

	var meta = make(map[string]string, 3)
	if invoke.opts.version != "" {
		meta[namespace.DefaultHeaderVersion] = invoke.opts.version
	}
	if invoke.opts.address != "" {
		meta[namespace.DefaultHeaderAddress] = invoke.opts.address
	}

	meta[namespace.DefaultHeaderContentType] = c.ContentType

	// clone context metadata
	cc, err := c.WithMetadata(invoke.opts.serviceName, method, meta)
	return cc, invoke, err
}

func (invoke *Invoke) Opts() InvokeOption {
	return invoke.opts
}

func (invoke *Invoke) Address() string {
	return invoke.opts.address
}

func (invoke *Invoke) Id() string {
	return invoke.opts.id
}

func (invoke *Invoke) Scope() string {
	return invoke.opts.scope
}

func (invoke *Invoke) FF() bool {
	return invoke.opts.ff
}

// InvokeRR invokeRR is invokeRequestResponse
func (invoke *Invoke) InvokeRR(c *context.Context, req, rsp proto.Message, cnn *conn.Conn) error {
	// encoding req body to roc packet
	b, err := c.Codec().Encode(req)
	if err != nil {
		c.Error(err)
		return err
	}
	var request, response = parcel.Payload(b), parcel.NewPacket()

	err = invokeRR(c, cnn, invoke, request, response, rsp)

	parcel.Recycle(request)
	parcel.Recycle(response)

	return err
}

func invokeRR(
	c *context.Context,
	cnn *conn.Conn,
	invoke *Invoke,
	request, response *parcel.RocPacket,
	rsp proto.Message,
) error {

	// send a request by requestResponse
	err := cnn.Client().RR(c, request, response)
	if err != nil {
		if invoke.opts.retry > 0 {
			// to retry request with backoff
			bf := backoff.NewBackoff()
			for i := 0; i < invoke.opts.retry; i++ {
				time.Sleep(bf.Next(i))
				if err = cnn.Client().RR(c, request, response); err == nil {
					break
				}
			}

			if err != nil {
				c.Error(err)

				// mark error count to manager conn state
				cnn.GrowError()
				return err
			}
		}
		return err
	}

	return c.Codec().Decode(response.Bytes(), rsp)
}

// InvokeFF invokeFF is FireAndForget
func (invoke *Invoke) InvokeFF(c *context.Context, req proto.Message, cnn *conn.Conn) {
	// encoding req body to roc packet
	b, err := c.Codec().Encode(req)
	if err != nil {
		c.Error(err)
		return
	}
	var request = parcel.Payload(b)

	// send a request by FireAndForget
	cnn.Client().FF(c, request)

	// defer recycle packet to pool
	parcel.Recycle(request)
}
