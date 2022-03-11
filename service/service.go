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

package service

import (
    "os"
    "os/signal"

    "github.com/go-roc/roc/config"
    "github.com/go-roc/roc/internal/etcd"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/rlog/log"
    "github.com/go-roc/roc/service/client"
    "github.com/go-roc/roc/service/opt"
    "github.com/go-roc/roc/service/server"
)

type Service struct {
    //service options setting
    opts opt.Option

    //service exit channel
    exit chan struct{}

    //roc service client,for rpc to server
    client *client.Client

    //roc service server,listen and wait call
    server *server.Server
}

func New(opts ...opt.Options) *Service {
    s := &Service{
        opts: opt.NewOpts(opts...),
        exit: make(chan struct{}),
    }

    s.server = server.NewServer(s.opts)
    return s
}

func (s *Service) Client() *client.Client {
    if s.client == nil {
        s.client = client.NewClient(s.opts)
    }

    return s.client
}

func (s *Service) Server() *server.Server {
    return s.server
}

func (s *Service) Run() {
    defer func() {
        if r := recover(); r != nil {
            rlog.Stack(r)
        }
    }()

    // handler signal
    ch := make(chan os.Signal)
    signal.Notify(ch, s.opts.Signal...)

    go func() {
        select {
        case c := <-ch:

            rlog.Infof("received signal %s ,service [%s] exit!", c.String(), s.opts.Name)

            s.CloseService()

            for _, f := range s.opts.Exit {
                f()
            }

            s.exit <- struct{}{}
        }
    }()

    //run server
    s.server.Run()

    select {
    case <-s.exit:
    }

    os.Exit(0)
}

func (s *Service) CloseService() {
    //close registry service discover
    if s.opts.Registry != nil {
        _ = s.opts.Registry.Deregister(s.opts.Endpoint)
        s.opts.Registry.CloseRegistry()
        s.opts.Registry = nil
    }

    //close service client
    if s.client != nil {
        s.client.CloseClient()
    }

    //close service server
    if s.server != nil {
        s.server.CloseServer()
    }

    //close config setting
    config.Close()

    //close etcd client
    etcd.DefaultEtcd.CloseEtcd()

    //todo flush rlog content
    log.Close()
}
