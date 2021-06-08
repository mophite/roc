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
	scope           string
	address         string
	service         string
	version         string
	buffSize        int
}

func BuffSize(buffSize int) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		invokeOption.buffSize = buffSize
	}
}

func WithScope(name string, version ...string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		var ver = namespace.DefaultVersion
		if len(version) == 1 {
			ver = version[0]
		}

		invokeOption.scope = namespace.SplicingScope(name, ver)
		invokeOption.service = name
		invokeOption.version = ver
	}
}

func WithAddress(name, address string, version ...string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		var ver = namespace.DefaultVersion
		if len(version) == 1 {
			ver = version[0]
		}

		invokeOption.scope = namespace.SplicingScope(name, ver)
		invokeOption.address = address
		invokeOption.service = name
		invokeOption.version = ver
	}
}

type Options func(*option)

type option struct {
	connectTimeout    time.Duration
	keepaliveInterval time.Duration
	keepaliveLifetime time.Duration
	client            transport.Client
	registry          registry.Registry
	retry             int
	cc                codec.Codec
}

func Client(client transport.Client) Options {
	return func(option *option) {
		option.client = client
	}
}

func Registry(registry registry.Registry) Options {
	return func(option *option) {
		option.registry = registry
	}
}

func Transport(client transport.Client) Options {
	return func(option *option) {
		option.client = client
	}
}

func ConnectTimeout(connectTimeout time.Duration) Options {
	return func(option *option) {
		option.connectTimeout = connectTimeout
	}
}

func Codec(cc codec.Codec) Options {
	return func(option *option) {
		option.cc = cc
	}
}

func KeepaliveInterval(keepaliveInterval time.Duration) Options {
	return func(option *option) {
		option.keepaliveInterval = keepaliveInterval
	}
}

func KeepaliveLifetime(keepaliveLifetime time.Duration) Options {
	return func(option *option) {
		option.keepaliveLifetime = keepaliveLifetime
	}
}

func newOpts(opts ...Options) option {
	opt := option{}
	for i := range opts {
		opts[i](&opt)
	}

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
		opt.client = rs.NewClient(
			opt.connectTimeout,
			opt.keepaliveInterval,
			opt.keepaliveLifetime,
		)
	}

	if opt.registry == nil {
		opt.registry = registry.DefaultEtcdRegistry(registry.Schema(namespace.DefaultSchema))
	}

	if opt.cc == nil {
		opt.cc = codec.DefaultCodec
	}

	return opt
}
