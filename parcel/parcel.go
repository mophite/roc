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

package parcel

import (
	"github.com/gogo/protobuf/proto"

	"roc/parcel/context"
)

type Handler func(c *context.Context, req *RocPacket, interrupt Interceptor) (rsp proto.Message, err error)

type StreamHandler func(c *context.Context, req *RocPacket) (chan proto.Message, chan error)

type ChannelHandler func(c *context.Context, req chan *RocPacket, errs chan error) (chan proto.Message, chan error)

type Fire func(c *context.Context, req proto.Message) (proto.Message, error)

type Interceptor func(c *context.Context, req proto.Message, fire Fire) (proto.Message, error)

type Wrapper func(c *context.Context) error
