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

package log

import (
	"github.com/go-roc/roc/rlog/format"
	"github.com/go-roc/roc/rlog/output"
)

type Option struct {
	call   int
	name   string
	prefix string
	format format.Formatter
	out    output.Outputor
}

type Options func(*Option)

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.format == nil {
		opt.format = format.DefaultFormat
	}

	if opt.name == "" {
		opt.name = "roc"
	}

	if opt.out == nil {
		opt.out = output.DefaultOutput
	}

	if opt.prefix == "" {
		opt.prefix = ""
	}

	if opt.call == 0 {
		opt.call = 4
	}

	return opt
}

func Call(call int) Options {
	return func(option *Option) {
		option.call = call
	}
}

func Output(out output.Outputor) Options {
	return func(option *Option) {
		option.out = out
	}
}

func Name(name string) Options {
	return func(option *Option) {
		option.name = name
	}
}

func Prefix(prefix string) Options {
	return func(option *Option) {
		option.prefix = prefix
	}
}

func Format(format format.Formatter) Options {
	return func(option *Option) {
		option.format = format
	}
}
