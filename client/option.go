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
	"time"

	"roc/internal/namespace"
	"roc/internal/registry"
	"roc/internal/transport"
	rs "roc/internal/transport/rscoket"
	"roc/parcel/codec"
)

type InvokeOptions func(*InvokeOption)

type InvokeOption struct {

	//scope is the service discovery prefix key
	scope string

	//address is witch server you want to call
	address string

	//serviceName is witch server by service serviceName
	serviceName string

	//version is witch server by version
	version string

	//buffSize effective only requestChannel
	buffSize int

	//trace is request flow trace
	//it's will be from web client,or generated on initialize
	trace string
}

// WithTracing set tracing
func WithTracing(t string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		invokeOption.trace = t
	}
}

// BuffSize set buff size for requestChannel
func BuffSize(buffSize int) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		invokeOption.buffSize = buffSize
	}
}

// WithName set service discover prefix with service serviceName
func WithName(name string, version ...string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		var ver = namespace.DefaultVersion

		// if no version ,use default version number
		if len(version) == 1 {
			ver = version[0]
		}

		invokeOption.scope = namespace.SplicingScope(name, ver)
		invokeOption.serviceName = name
		invokeOption.version = ver
	}
}

// WithAddress set service discover prefix with both service serviceName and address
func WithAddress(name, address string, version ...string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		var ver = namespace.DefaultVersion

		// if no version ,use default version number
		if len(version) == 1 {
			ver = version[0]
		}

		invokeOption.scope = namespace.SplicingScope(name, ver)
		invokeOption.address = address
		invokeOption.serviceName = name
		invokeOption.version = ver
	}
}

// Options invoke option
type Options func(*option)

// invoke option
type option struct {

	// connect server within connectTimeout
	// if out of ranges,will be timeout
	connectTimeout time.Duration

	// keepalive setting,the period for requesting heartbeat to stay connected
	keepaliveInterval time.Duration

	// keepalive setting,the longest time the connection can survive
	keepaliveLifetime time.Duration

	// transport client
	client transport.Client

	//service discover registry
	registry registry.Registry

	//for requestResponse try to retry request
	retry int

	//data encoding or decoding
	cc codec.Codec
}

// Registry set service discover registry
func Registry(registry registry.Registry) Options {
	return func(option *option) {
		option.registry = registry
	}
}

// Transport set transport client
func Transport(client transport.Client) Options {
	return func(option *option) {
		option.client = client
	}
}

// ConnectTimeout set connect timeout
func ConnectTimeout(connectTimeout time.Duration) Options {
	return func(option *option) {
		option.connectTimeout = connectTimeout
	}
}

// Codec set codec
func Codec(cc codec.Codec) Options {
	return func(option *option) {
		option.cc = cc
	}
}

// KeepaliveInterval set keepalive interval
func KeepaliveInterval(keepaliveInterval time.Duration) Options {
	return func(option *option) {
		option.keepaliveInterval = keepaliveInterval
	}
}

// KeepaliveLifetime set keepalive life time
func KeepaliveLifetime(keepaliveLifetime time.Duration) Options {
	return func(option *option) {
		option.keepaliveLifetime = keepaliveLifetime
	}
}

// new invoke option
func newOpts(opts ...Options) option {
	opt := option{}
	for i := range opts {
		opts[i](&opt)
	}

	//set connect timeout or default
	if opt.connectTimeout <= 0 {
		opt.connectTimeout = time.Second * 5
	}

	if opt.keepaliveLifetime <= 0 {
		opt.keepaliveLifetime = time.Second * 600
	}

	if opt.keepaliveInterval <= 0 {
		opt.keepaliveInterval = time.Second * 5
	}

	if opt.client == nil {
		//default is rsocket
		opt.client = rs.NewClient(
			opt.connectTimeout,
			opt.keepaliveInterval,
			opt.keepaliveLifetime,
		)
	}

	if opt.registry == nil {
		//default registry with default schema
		opt.registry = registry.NewRegistry(registry.Schema(namespace.DefaultSchema))
	}

	if opt.cc == nil {
		opt.cc = codec.DefaultCodec
	}

	return opt
}
