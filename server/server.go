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

package server

import (
	"errors"
	"os"
	"os/signal"
	"roc/parcel/codec"

	"roc/internal/endpoint"
	"roc/internal/router"
	"roc/parcel"

	"roc/internal/registry"
	"roc/internal/transport"
	"roc/rlog"
	"roc/rlog/log"
)

type RocServer struct {

	//run server option
	opts option

	//server exit channel
	exit chan struct{}

	//server router collection
	route *router.Router
}

func (r *RocServer) Id() string {
	return r.opts.id
}

func (r *RocServer) Name() string {
	return r.opts.name
}

func (r *RocServer) TCPAddress() string {
	return r.opts.tcpAddress
}

func (r *RocServer) Version() string {
	return r.opts.version
}

func (r *RocServer) Register() registry.Registry {
	return r.opts.register
}

func (r *RocServer) WssAddress() string {
	return r.opts.wssAddress
}

func (r *RocServer) WssPath() string {
	return r.opts.wssPath
}

func (r *RocServer) Transport() transport.Server {
	return r.opts.transporter
}

func (r *RocServer) E() *endpoint.Endpoint {
	return r.opts.e
}

func (r *RocServer) Codec() codec.Codec {
	return r.opts.cc
}

func NewRocServer(opts ...Options) *RocServer {
	r := &RocServer{
		opts: newOpts(opts...),
		exit: make(chan struct{}),
	}
	r.route = router.NewRouter(r.opts.wrappers, r.opts.err, r.opts.cc)

	//NOTICE: don't register wss to sd.
	//if r.SetupWss() {
	//	err := r.opts.register.Register(r.WssAddress(), "wss")
	//	if err != nil {
	//		return nil
	//	}
	//}
	r.opts.transporter.Accept(r.route)

	return r
}

func (r *RocServer) Run() error {
	defer func() {
		if r := recover(); r != nil {
			rlog.Stack(r)
		}
	}()

	// handler signal
	ch := make(chan os.Signal)
	signal.Notify(ch, r.opts.signal...)

	go func() {
		select {
		case c := <-ch:
			r.Close()

			for _, f := range r.opts.exit {
				f()
			}

			rlog.Infof("received signal %s [%s] server exit!", c.String(), r.opts.name)

			log.Close()

			r.exit <- struct{}{}
		}
	}()

	// echo method list
	r.route.List()
	r.opts.transporter.Run()

	rlog.Infof("[tcp:%s] AND [ws:%s] is start success!",
		r.opts.e.Absolute,
		r.opts.wssAddress,
	)

	err := r.register()
	if err != nil {
		panic(err)
	}

	select { case <-r.exit: }

	return errors.New(r.opts.name + " server is exit!")
}

func (r *RocServer) register() error {
	return r.Register().Register(r.opts.e)
}

func (r *RocServer) RegisterHandler(method string, rr parcel.Handler) {
	r.route.RegisterHandler(method, rr)
}

func (r *RocServer) RegisterStreamHandler(method string, rs parcel.StreamHandler) {
	r.route.RegisterStreamHandler(method, rs)
}

func (r *RocServer) RegisterChannelHandler(method string, rs parcel.ChannelHandler) {
	r.route.RegisterChannelHandler(method, rs)
}

func (r *RocServer) Close() {
	_ = r.Register().Deregister(r.opts.e)
	r.opts.register.Close()
	r.opts.transporter.Close()
}
