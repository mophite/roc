package client

import (
	"time"

	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/service/conn"
)

type Option struct {

	//data codec
	cc codec.Codec

	//discover service
	registry registry.Registry
}

type Options func(option *Option)

func Codec(cc codec.Codec) Options {
	return func(option *Option) {
		option.cc = cc
	}
}

func Registry(registry registry.Registry) Options {
	return func(option *Option) {
		option.registry = registry
	}
}

func ConnectTimeout(timeout time.Duration) Options {
	return func(option *Option) {
		conn.DefaultConnectTimeout = timeout
	}
}

func KeepaliveInterval(keepaliveInterval time.Duration) Options {
	return func(option *Option) {
		conn.DefaultKeepaliveInterval = keepaliveInterval
	}
}

func KeepaliveLifetime(keepaliveLifetime time.Duration) Options {
	return func(option *Option) {
		conn.DefaultKeepaliveLifetime = keepaliveLifetime
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	//if opt.ratelimit <= 0 {
	//	opt.ratelimit = math.MaxInt32
	//}
	//set connect timeout or default
	if opt.cc == nil {
		opt.cc = codec.DefaultCodec
	}

	if opt.registry == nil {
		opt.registry = registry.DefaultRegistry
	}

	return opt
}
