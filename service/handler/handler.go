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

package handler

import (
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/gogo/protobuf/proto"
)

//Handler for rpc service handler
type Handler func(c *context.Context, req *parcel.RocPacket, interrupt Interceptor) (rsp proto.Message, err error)

// StreamHandler for rpc service stream handler
type StreamHandler func(c *context.Context, req *parcel.RocPacket) (chan proto.Message, chan error)

// ChannelHandler for rpc service channel handler
type ChannelHandler func(c *context.Context, req chan *parcel.RocPacket, errs chan error) (chan proto.Message, chan error)

//Fire run interceptor action
type Fire func(c *context.Context, req proto.Message) (proto.Message, error)

// Interceptor for rpc request response interceptor function
type Interceptor func(c *context.Context, req proto.Message, fire Fire) (proto.Message, error)

// WrapperHandler for all rpc function middleware
type WrapperHandler func(c *context.Context) error

// ApiHandler register router by POST http method
//the request body with content-type
//will be json or proto data protocol. eg.
//func (h *hello)Say(c *Context,req *phello.SayReq,rsp *phello.SayRsp)error
type ApiHandler func(c *context.Context, req proto.Message, rsp proto.Message)

type ApiRocHandler func(c *context.Context) (rsp proto.Message, err error)