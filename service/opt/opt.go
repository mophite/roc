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
    "net/http"
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
    "github.com/rs/cors"
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
    //just for request response
    Wrappers []handler.WrapperHandler

    //just for http request before or socket setup before
    Dog []handler.DogHandler

    //when transportServer exit,will do exit func
    Exit []func()

    //etcd config
    EtcdConfig *clientv3.Config

    //config options
    ConfigOpt []config.Options

    //service discover registry
    Registry registry.Registry

    //service discover endpoint
    Endpoint *endpoint.Endpoint

    //only need cors middleware on roc http api POST/DELETE/GET/PUT/OPTIONS method
    CorsOptions *cors.Options
}

func NewOpts(opts ...Options) Option {
    opt := Option{}

    for i := range opts {
        opts[i](&opt)
    }

    err := config.NewConfig(opt.ConfigOpt...)
    if err != nil {
        panic("config NewConfig occur error: " + err.Error())
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
        opt.RandPort = &[2]int{10000, 59999}
    }

    // NOTICE: api service only support fixed tcpAddress ,not suggest rand tcpAddress in api service
    if opt.TcpAddress == "" {
        opt.TcpAddress = ip + ":" + strconv.Itoa(x.RandInt(opt.RandPort[0], opt.RandPort[1]))
    }

    if opt.WssAddress != "" && opt.WssPath == "" {
        opt.WssPath = "/roc/wss"
    }

    if opt.WssPath != "" {
        if !strings.HasPrefix(opt.WssPath, "/") {
            opt.WssPath = "/" + opt.WssPath
        }

        if strings.HasSuffix(opt.WssPath, "/") {
            opt.WssPath = strings.TrimSuffix(opt.WssPath, "/")
        }
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

    if opt.Registry == nil {
        opt.Registry = registry.NewRegistry()
    }

    if opt.CorsOptions == nil {
        //allowed all
        opt.CorsOptions = &cors.Options{
            AllowedOrigins: []string{"*"},
            AllowedMethods: []string{
                http.MethodGet,
                http.MethodPost,
                http.MethodPut,
                http.MethodDelete,
                http.MethodOptions,
            },
            AllowedHeaders:   []string{"*"},
            AllowCredentials: false,
            Debug:            false,
        }
    }

    return opt
}
