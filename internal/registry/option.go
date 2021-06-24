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

	"github.com/go-roc/roc/internal/endpoint"
)

type Option struct {

	//endpoint
	e *endpoint.Endpoint

	//timeout setting
	timeout time.Duration

	//schema prefix
	schema string
}

type Options func(option *Option)

func Timeout(timeout time.Duration) Options {
	return func(option *Option) {
		option.timeout = timeout
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

	if opt.schema == "" {
		panic("no schema setting")
	}

	return opt
}
