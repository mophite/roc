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

package roc

import (
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/go-roc/roc/internal/backoff"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/context"
)

type Invoke struct {

	// strategy clients to invoke
	strategy Strategy

	// invoke options
	opts InvokeOption
}

// create a invoke
func newInvoke(c *context.Context, method string, service *Service, opts ...InvokeOptions) (*Invoke, error) {
	invoke := &Invoke{strategy: service.strategy}

	for i := range opts {
		opts[i](&invoke.opts)
	}

	// initialize tunnel for requestChannel only
	if invoke.opts.buffSize == 0 {
		invoke.opts.buffSize = 10
	}

	// create metadata
	var err = c.WithMetadata(
		invoke.opts.serviceName,
		method,
		invoke.opts.trace,
		map[string]string{
			namespace.DefaultHeaderVersion: invoke.opts.version,
			namespace.DefaultHeaderAddress: invoke.opts.address,
		},
	)
	return invoke, err
}

// invokeRR is invokeRequestResponse
func (i *Invoke) invokeRR(c *context.Context, req, rsp proto.Message, conn *Conn, opts Option) error {
	// encoding req body to roc packet
	b, err := opts.cc.Encode(req)
	if err != nil {
		return err
	}

	var request, response = parcel.Payload(b), parcel.NewPacket()

	// defer recycle packet to pool
	defer func() {
		parcel.Recycle(response, request)
	}()

	// send a request by requestResponse
	err = conn.Client().RR(c, request, response)
	if err != nil {

		// to retry request with backoff
		bf := backoff.NewBackoff()
		for i := 0; i < opts.retry; i++ {
			time.Sleep(bf.Next(i))
			if err = conn.Client().RR(c, request, response); err == nil {
				break
			}
		}

		if err != nil {
			c.Error(err)

			// mark error count to manager conn state
			conn.growError()
			return err
		}
	}

	return opts.cc.Decode(response.Bytes(), rsp)
}
