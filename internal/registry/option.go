package registry

import (
	"time"

	"go.etcd.io/etcd/clientv3"

	"roc/internal/endpoint"
	"roc/internal/namespace"
)

type Option struct {
	e          *endpoint.Endpoint
	version    string
	name       string
	address    []string
	timeout    time.Duration
	leaseTLL   int64
	schema     string
	etcdConfig *clientv3.Config
}

type Options func(option *Option)

func EtcdConfig(c *clientv3.Config) Options {
	return func(option *Option) {
		option.etcdConfig = c
	}
}

func Name(name string) Options {
	return func(option *Option) {
		option.name = name
	}
}

func Version(version string) Options {
	return func(option *Option) {
		option.version = version
	}
}

func Address(address []string) Options {
	return func(option *Option) {
		option.address = address
	}
}

func Timeout(timeout time.Duration) Options {
	return func(option *Option) {
		option.timeout = timeout
	}
}

func LeaseTLL(leaseTLL int64) Options {
	return func(option *Option) {
		option.leaseTLL = leaseTLL
	}
}

func Schema(schema string) Options {
	return func(option *Option) {
		option.schema = schema
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.timeout == 0 {
		opt.timeout = time.Second * 5
	}

	if opt.leaseTLL == 0 {
		opt.leaseTLL = 5
	}

	if opt.version == "" {
		opt.version = namespace.DefaultVersion
	}

	if opt.schema == "" {
		opt.schema = namespace.DefaultSchema
	}

	return opt
}
