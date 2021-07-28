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
    "github.com/go-roc/roc/parcel/codec"

    "github.com/go-roc/roc/parcel/packet"
)

type ErrorCode = int32

const (
    ErrorCodeInternalServer ErrorCode = 500
    ErrorCodeBadRequest               = 400
    ErrCodeNotFoundHandler            = 401
)

var (
    DefaultErrorPacket = NewErrorPacket()
)

type ErrorPackager interface {
    Encode(c codec.Codec, code ErrorCode, err error) []byte
}

type ErrorPacket struct {
    data packet.Packet
}

func NewErrorPacket() *ErrorPacket {
    return &ErrorPacket{}
}

func (e *ErrorPacket) Encode(c codec.Codec, code ErrorCode, err error) []byte {
    ep := *e
    ep.data.Code = code
    ep.data.Msg = err.Error()
    b, _ := c.Encode(&ep.data)
    return b
}
