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

package roc

import (
    "os"
    "os/signal"
    "time"

    "github.com/coreos/etcd/clientv3"
    "github.com/go-roc/roc/config"
    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/etcd"
    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/internal/registry"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/codec"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/rlog/log"
    "github.com/go-roc/roc/service/client"
    "github.com/go-roc/roc/service/conn"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/opt"
    "github.com/go-roc/roc/service/server"
    "github.com/rs/cors"
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

const SupportPackageIsVersion1 = 1

func Release() {
    log.SetInfo()
}

func BuffSize(buffSize int) opt.Options {
    return func(o *opt.Option) {
        o.BuffSize = buffSize
    }
}

func Wrapper(wrappers ...handler.WrapperHandler) opt.Options {
    return func(o *opt.Option) {
        o.Wrappers = append(o.Wrappers, wrappers...)
    }
}

func WatchDog(wrappers ...handler.DogHandler) opt.Options {
    return func(o *opt.Option) {
        o.Dog = wrappers
    }
}

func Exit(exit ...func()) opt.Options {
    return func(o *opt.Option) {
        o.Exit = exit
    }
}

func Signal(signal ...os.Signal) opt.Options {
    return func(o *opt.Option) {
        o.Signal = signal
    }
}

// Port port[0]:min port[1]:max
func Port(port [2]int) opt.Options {
    return func(o *opt.Option) {
        if port[0] > port[1] {
            panic("port[1] must greater than port[0]")
        }

        if port[0] < 10000 {
            panic("rand port for internal transportServer suggest more than 10000")
        }

        o.RandPort = &port
    }
}

func Error(err parcel.ErrorPackager) opt.Options {
    return func(o *opt.Option) {
        o.Err = err
    }
}

func WssApiAddr(address, path string) opt.Options {
    return func(o *opt.Option) {
        o.WssAddress = address
        o.WssPath = path
    }
}

func HttpApiAddr(address string) opt.Options {
    return func(o *opt.Option) {
        o.HttpAddress = address
    }
}

func Id(id string) opt.Options {
    return func(o *opt.Option) {
        o.Id = id
    }
}

func Namespace(name string) opt.Options {
    return func(o *opt.Option) {
        o.Name = name
    }
}

func RootRouterRedirect(r string) opt.Options {
    return func(option *opt.Option) {
        option.RootRouterRedirect = r
    }
}

// EtcdConfig setting global etcd config first
func EtcdConfig(e *clientv3.Config) opt.Options {
    return func(o *opt.Option) {
        o.EtcdConfig = e
    }
}

func TCPApiSrvPort(port int) opt.Options {
    return func(o *opt.Option) {
        o.TcpAddress = port
    }
}

func ConfigOption(opts ...config.Options) opt.Options {
    return func(o *opt.Option) {
        o.ConfigOpt = opts
    }
}

func Codec(contentType string, c codec.Codec) opt.Options {
    return func(o *opt.Option) {
        codec.SetCodec(contentType, c)
    }
}

func Version(version string) opt.Options {
    return func(o *opt.Option) {
        namespace.DefaultVersion = version
    }
}

func Registry(r registry.Registry) opt.Options {
    return func(o *opt.Option) {
        o.Registry = r
    }
}

// Deprecated: newServer will created endpoint
func Endpoint(e *endpoint.Endpoint) opt.Options {
    return func(o *opt.Option) {
        o.Endpoint = e
    }
}

func ConnectTimeout(timeout time.Duration) opt.Options {
    return func(o *opt.Option) {
        conn.DefaultConnectTimeout = timeout
    }
}

func KeepaliveInterval(keepaliveInterval time.Duration) opt.Options {
    return func(o *opt.Option) {
        conn.DefaultKeepaliveInterval = keepaliveInterval
    }
}

func KeepaliveLifetime(keepaliveLifetime time.Duration) opt.Options {
    return func(o *opt.Option) {
        conn.DefaultKeepaliveLifetime = keepaliveLifetime
    }
}

func Cors(m *cors.Options) opt.Options {
    return func(option *opt.Option) {
        option.CorsOptions = m
    }
}
