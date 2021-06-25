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

package transport

import (
	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/router"
	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/context"
)

type Server interface {

	// Address about server socket setup on ip address
	Address() string

	Accept(fn *router.Router)

	Run()

	String() string

	Close()
}

type Client interface {
	Dial(e *endpoint.Endpoint, closeChan chan string) error

	// RR request/response,through block unsafe method
	RR(c *context.Context, req *parcel.RocPacket, rsp *parcel.RocPacket) (err error)

	// FF FireAndForget
	//FF(c *context.Context, req *parcel.RocPacket)

	// RS request/stream
	RS(c *context.Context, req *parcel.RocPacket) (chan []byte, chan error)

	// RC request/channel
	RC(c *context.Context, req chan []byte, errsIn chan error) (chan []byte, chan error)

	// MP metadata
	//MP(c *context.Context)

	String() string

	Close()
}

// CallOptions todo call server Options
type CallOptions func(option *CallOption)

type CallOption struct {
}
