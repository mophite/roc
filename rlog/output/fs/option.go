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

package fs

import (
	"os"
	"path/filepath"
	"time"
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
)

const (
	DefaultFileMaxDelta = 100
	DefaultBufferSize   = 512 * KB
	DefaultFileMaxSize  = 256 * MB
)

type Option struct {
	link          string
	name          string // project name
	fileName      string
	path          string
	async         bool
	prefix        string
	maxFileSize   int
	maxBufferSize int
	maxBucketSize int
	rotate        bool
	interval      time.Duration
	zone          *time.Location
	modePerm      int
}

type Options func(*Option)

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.name == "" {
		opt.name = "roc"
	}

	if opt.path == "" {
		opt.path = "./logs"
	}

	opt.link = opt.path + string(os.PathSeparator) + opt.name + ".log"

	if opt.maxFileSize == 0 {
		opt.maxFileSize = DefaultFileMaxSize
	}

	opt.maxFileSize -= DefaultFileMaxDelta

	if opt.maxBufferSize == 0 {
		opt.maxBufferSize = DefaultBufferSize
	}

	if opt.interval == 0 {
		opt.interval = time.Millisecond * 500
	}

	if opt.zone == nil {
		opt.zone = time.Local
	}

	if opt.prefix == "" {
		opt.prefix = ""
	}

	if opt.modePerm == 0 {
		opt.modePerm = int(os.ModePerm)
	}

	opt.link = filepath.Join(opt.path, opt.name+".log")
	return opt
}

func Name(name string) Options {
	return func(option *Option) {
		option.name = name
	}
}

func Interval(interval time.Duration) Options {
	return func(option *Option) {
		option.interval = interval
	}
}

func Link(link string) Options {
	return func(option *Option) {
		option.link = link
	}
}

func Path(p string) Options {
	return func(option *Option) {
		option.path = p
	}
}

func Async() Options {
	return func(option *Option) {
		option.async = true
	}
}

func Prefix(prefix string) Options {
	return func(option *Option) {
		option.prefix = prefix
	}
}

func MaxFileSize(maxFileSize int) Options {
	return func(option *Option) {
		option.maxBufferSize = maxFileSize
	}
}

func MaxBufferSize(maxBufferSize int) Options {
	return func(option *Option) {
		option.maxBufferSize = maxBufferSize
	}
}

func Rotate() Options {
	return func(option *Option) {
		option.rotate = true
	}
}

func Zone(zone *time.Location) Options {
	return func(option *Option) {
		option.zone = zone
	}
}

func FileModel(perm int) Options {
	return func(option *Option) {
		option.modePerm = perm
	}
}
