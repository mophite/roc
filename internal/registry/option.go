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

package registry

import (
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/namespace"
)

type Option struct {

	//endpoint
	e *endpoint.Endpoint

	//version prefix
	version string

	//name prefix
	name string

	//address prefix
	address []string

	//timeout setting
	timeout time.Duration

	//lease setting
	leaseTLL int64

	//schema prefix
	schema string

	//clientv3 confi
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
