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


package opt

import (
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/coreos/etcd/clientv3"
    "github.com/go-roc/roc/config"
    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/etcd"
    "github.com/go-roc/roc/internal/registry"
    "github.com/go-roc/roc/internal/sig"
    "github.com/go-roc/roc/internal/transport"
    rs "github.com/go-roc/roc/internal/transport/rsocket"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/x"
    "github.com/go-roc/roc/x/fs"
    "github.com/gorilla/mux"
)

type Options func(option *Option)

type Option struct {
    //transportServer id can be set,or random
    //this is the unique identifier of the service
    Id string

    //transportServer name,eg.srv.hello or api.hello
    Name string

    //random port
    //Random ports make you don’t need to care about specific port numbers,
    //which are more commonly used in internal services
    RandPort *[2]int

    //socket tcp ip:port address
    TcpAddress string

    //websocket ip:port address
    WssAddress string

    //websocket relative path address of websocket
    WssPath string

    //http service address
    HttpAddress string

    //buffSize to data tunnel if it's need
    BuffSize int

    //transportServer transport
    TransportServer transport.Server

    //error packet
    //It will affect the format of the data you return
    Err parcel.ErrorPackager

    //receive system signal
    Signal []os.Signal

    //wrapper some middleware
    //it's can be interrupt
    Wrappers []handler.WrapperHandler

    //when transportServer exit,will do exit func
    Exit []func()

    //it must be unique in all of your handler path
    ApiPrefix string

    //etcd config
    EtcdConfig *clientv3.Config

    //config options
    ConfigOpt []config.Options

    //http api router
    Router *mux.Router

    //service discover registry
    Registry registry.Registry

    //service discover endpoint
    Endpoint *endpoint.Endpoint
}

func NewOpts(opts ...Options) Option {
    opt := Option{}

    for i := range opts {
        opts[i](&opt)
    }

    if opt.Name == "" {
        opt.Name = fs.GetProjectName()
    }

    if opt.Id == "" {
        //todo change to git commit id+timestamp
        opt.Id = x.NewUUID()
    }

    ip, err := x.LocalIp()
    if err != nil {
        panic(err)
    }

    if opt.RandPort == nil {
        opt.RandPort = &[2]int{10000, 19999}
    }

    // NOTICE: api service only support fixed tcpAddress ,not suggest rand tcpAddress in api service
    if opt.TcpAddress == "" {
        opt.TcpAddress = ip + ":" + strconv.Itoa(x.RandInt(opt.RandPort[0], opt.RandPort[1]))
    }

    if opt.WssAddress != "" && opt.WssPath == "" {
        opt.WssPath = "/roc/wss"
    }

    if opt.Err == nil {
        opt.Err = parcel.DefaultErrorPacket
    }

    if opt.BuffSize == 0 {
        opt.BuffSize = 10
    }

    if opt.TransportServer == nil {
        opt.TransportServer = rs.NewServer(opt.TcpAddress, opt.WssAddress, opt.Name, opt.BuffSize)
    }

    if opt.Endpoint == nil {
        opt.Endpoint, err = endpoint.NewEndpoint(opt.Id, opt.Name, opt.TcpAddress)
        if err != nil {
            panic(err)
        }
    }

    if opt.Signal == nil {
        opt.Signal = sig.DefaultSignal
    }

    if opt.ApiPrefix == "" {
        opt.ApiPrefix = "roc"
    }

    if !strings.HasPrefix(opt.ApiPrefix, "/") {
        opt.ApiPrefix = "/" + opt.ApiPrefix
    }

    if !strings.HasSuffix(opt.ApiPrefix, "/") {
        opt.ApiPrefix = opt.ApiPrefix + "/"
    }

    if opt.EtcdConfig == nil {
        opt.EtcdConfig = &clientv3.Config{
            Endpoints:   []string{"127.0.0.1:2379"},
            DialTimeout: time.Second * 5,
        }
    }

    // init etcd.DefaultEtcd
    err = etcd.NewEtcd(time.Second*5, 5, opt.EtcdConfig)
    if err != nil {
        panic("etcdConfig occur error: " + err.Error())
    }

    err = config.NewConfig(opt.ConfigOpt...)
    if err != nil {
        panic("config NewConfig occur error: " + err.Error())
    }

    opt.Router = mux.NewRouter()

    if opt.Registry == nil {
        opt.Registry = registry.NewRegistry()
    }

    return opt
}
