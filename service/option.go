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
    "strings"
    "time"

    "github.com/coreos/etcd/clientv3"
    "github.com/go-roc/roc/config"
    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/internal/registry"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/codec"
    "github.com/go-roc/roc/service/conn"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/opt"
)

const SupportPackageIsVersion1 = 1

//DefaultApiPrefix it must be unique in all of your handler path
var DefaultApiPrefix = "/roc/"

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

func Port(port [2]int) opt.Options {
    return func(o *opt.Option) {
        if port[0] > port[1] {
            panic("port index 0 must more than 1")
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

func WssAddress(address, path string) opt.Options {
    return func(o *opt.Option) {
        o.WssAddress = address
        o.WssPath = path
    }
}

func HttpAddress(address string) opt.Options {
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

// EtcdConfig setting global etcd config first
func EtcdConfig(e *clientv3.Config) opt.Options {
    return func(o *opt.Option) {
        o.EtcdConfig = e
    }
}

func TCPAddress(address string) opt.Options {
    return func(o *opt.Option) {
        o.TcpAddress = address
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

func ApiPrefix(apiPrefix string) opt.Options {
    return func(o *opt.Option) {
        if !strings.HasPrefix(apiPrefix, "/") {
            apiPrefix = "/" + apiPrefix
        }
        if !strings.HasSuffix(apiPrefix, "/") {
            apiPrefix += "/"
        }
        DefaultApiPrefix = apiPrefix
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
