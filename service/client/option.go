package client

import (
	"time"

	"github.com/go-roc/roc/config"
)

type Option struct {
	// connect server within connectTimeout
	// if out of ranges,will be timeout
	connectTimeout time.Duration

	// keepalive setting,the period for requesting heartbeat to stay connected
	keepaliveInterval time.Duration

	// keepalive setting,the longest time the connection can survive
	keepaliveLifetime time.Duration

	//config options
	configOpt []config.Options
}

type Options func(option *Option)

func ConfigOption(opts ...config.Options) Options {
	return func(option *Option) {
		option.configOpt = opts
	}
}

// ConnectTimeout set connect timeout
func ConnectTimeout(connectTimeout time.Duration) Options {
	return func(option *Option) {
		option.connectTimeout = connectTimeout
	}
}

// KeepaliveInterval set keepalive interval
func KeepaliveInterval(keepaliveInterval time.Duration) Options {
	return func(option *Option) {
		option.keepaliveInterval = keepaliveInterval
	}
}

// KeepaliveLifetime set keepalive life time
func KeepaliveLifetime(keepaliveLifetime time.Duration) Options {
	return func(option *Option) {
		option.keepaliveLifetime = keepaliveLifetime
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	err := config.NewConfig(opt.configOpt...)
	if err != nil {
		panic("config NewConfig occur error: " + err.Error())
	}

	//if opt.ratelimit <= 0 {
	//	opt.ratelimit = math.MaxInt32
	//}
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

	return opt
}
