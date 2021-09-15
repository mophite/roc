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

package invoke

import (
    "strings"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/parcel/codec"
)

type InvokeOptions func(*InvokeOption)

type InvokeOption struct {

    //scope is the service discovery prefix key
    scope string

    //address is witch server you want to call
    address string

    //serviceName is witch server by service serviceName
    serviceName string

    //version is witch server by version
    version string

    //buffSize effective only requestChannel
    buffSize int

    trace string

    prefix string

    //for requestResponse try to retry request
    retry int

    //data encoding or decoding
    cc codec.Codec
}

func Codec(cc codec.Codec) InvokeOptions {
    return func(option *InvokeOption) {
        option.cc = cc
    }
}

// WithTracing set tracing
func WithTracing(t string) InvokeOptions {
    return func(invokeOption *InvokeOption) {
        invokeOption.trace = t
    }
}

// InvokeBuffSize set buff size for requestChannel
func InvokeBuffSize(buffSize int) InvokeOptions {
    return func(invokeOption *InvokeOption) {
        invokeOption.buffSize = buffSize
    }
}

// WithName set service discover prefix with service serviceName
func WithName(name string, version ...string) InvokeOptions {
    return func(invokeOption *InvokeOption) {
        var ver = namespace.DefaultVersion

        // if no version ,use default version number
        if len(version) == 1 {
            ver = version[0]
        }

        invokeOption.scope = name + "/" + ver
        invokeOption.serviceName = name
        invokeOption.version = ver

        invokeOption.prefix = name

        ss := strings.Split(invokeOption.prefix, ".")
        invokeOption.prefix = ss[len(ss)-1]

        if strings.HasSuffix(invokeOption.prefix, "/") {
            invokeOption.prefix = strings.TrimSuffix(invokeOption.prefix, "/")
        }

        if !strings.HasPrefix(invokeOption.prefix, "/") {
            invokeOption.prefix = "/" + invokeOption.prefix
        }
    }
}

// WithAddress set service discover prefix with both service serviceName and address
func WithAddress(name, address string, version ...string) InvokeOptions {
    return func(invokeOption *InvokeOption) {
        var ver = namespace.DefaultVersion

        // if no version ,use default version number
        if len(version) == 1 {
            ver = version[0]
        }

        invokeOption.scope = name + "/" + ver
        invokeOption.address = address
        invokeOption.serviceName = name
        invokeOption.version = ver
    }
}
