package service

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/go-roc/roc/config"
	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/etcd"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/rlog/log"
	"github.com/go-roc/roc/service/client"
	"github.com/go-roc/roc/service/server"
)

type Option struct {

	//etcd config
	etcdConfig *clientv3.Config

	//config options
	configOpt []config.Options

	//it must be unique in all of your handler path
	apiPrefix string

	//service version
	version string

	e *endpoint.Endpoint
}

type Options func(option *Option)

// EtcdConfig setting global etcd config first
func EtcdConfig(e *clientv3.Config) Options {
	return func(option *Option) {
		option.etcdConfig = e
	}
}

func Codec(cc codec.Codec) Options {
	return func(option *Option) {
		codec.DefaultCodec = nil
		codec.DefaultCodec = cc
	}
}

func ConfigOption(opts ...config.Options) Options {
	return func(option *Option) {
		option.configOpt = opts
	}
}

func Version(version string) Options {
	return func(option *Option) {
		namespace.DefaultVersion = version
	}
}

func Registry(r registry.Registry) Options {
	return func(option *Option) {
		registry.DefaultRegistry = nil
		registry.DefaultRegistry = r
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.etcdConfig == nil {
		opt.etcdConfig = &clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Second * 5,
		}
	}

	// init etcd.DefaultEtcd
	err := etcd.NewEtcd(time.Second*5, 5, opt.etcdConfig)
	if err != nil {
		panic("etcdConfig occur error: " + err.Error())
	}

	err = config.NewConfig(opt.configOpt...)
	if err != nil {
		panic("config NewConfig occur error: " + err.Error())
	}

	return opt
}

type Service struct {
	opts Option

	client client.Client

	server server.Server
}

func (s *Service) Close() {
	if registry.DefaultRegistry != nil {
		_ = registry.DefaultRegistry.Deregister(s.opts.e)
		registry.DefaultRegistry.Close()
	}

	//todo flush rlog content
	log.Close()
}
