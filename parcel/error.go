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
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/x"

    "github.com/go-roc/roc/parcel/packet"
)

var (
    DefaultErrorPacket = NewErrorPacket()
)

type ErrorPackager interface {
    Error400(c *context.Context) []byte
    Error500(c *context.Context) []byte
    Error404(c *context.Context) []byte
    Error405(c *context.Context) []byte
}

type ErrorPacket struct{}

func (e *ErrorPacket) Error400(c *context.Context) []byte {
    p := new(packet.Packet)
    p.Code = 400
    p.Msg = "Bad Request"
    if c == nil {
        return x.MustMarshal(p)
    }
    return c.Codec().MustEncode(p)
}

func (e *ErrorPacket) Error500(c *context.Context) []byte {
    p := new(packet.Packet)
    p.Code = 500
    p.Msg = "Internal server error"
    if c == nil {
        return x.MustMarshal(p)
    }
    return c.Codec().MustEncode(p)
}

func (e *ErrorPacket) Error404(c *context.Context) []byte {
    p := new(packet.Packet)
    p.Code = 404
    p.Msg = "Not Found"
    if c == nil {
        return x.MustMarshal(p)
    }
    return c.Codec().MustEncode(p)
}

func (e *ErrorPacket) Error405(c *context.Context) []byte {
    p := new(packet.Packet)
    p.Code = 405
    p.Msg = "Method Not Allowed"
    if c == nil {
        return x.MustMarshal(p)
    }
    return c.Codec().MustEncode(p)
}

func NewErrorPacket() *ErrorPacket {
    return &ErrorPacket{}
}
