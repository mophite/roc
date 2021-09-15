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
    "bytes"
    "io"
)

const (
    //default bytes buffer cap size
    defaultCapSize = 512
    //default packet pool size
    defaultPoolSize = 1024000
    //default bytes buffer cap max size
    //if defaultCapSize > defaultMaxCapSize ? defaultCapSize
    defaultMaxCapSize = 4096
)

//create default pool
var pool = &packetPool{
    poolSize:   defaultPoolSize,
    capSize:    defaultCapSize,
    maxCapSize: defaultMaxCapSize,
    packets:    make(chan *RocPacket, defaultPoolSize),
}

type packetPool struct {
    poolSize, capSize, maxCapSize int
    packets                       chan *RocPacket
}

type RocPacket struct {
    B *bytes.Buffer
}

func newRocPacket() *RocPacket {
    return &RocPacket{B: bytes.NewBuffer(make([]byte, 0, pool.capSize))}
}

func Recycle(p ...*RocPacket) {
    for i := range p {
        p[i].B.Reset()

        if p[i].B.Cap() > pool.maxCapSize {
            p[i].B = bytes.NewBuffer(make([]byte, 0, pool.maxCapSize))
        }

        select {
        case pool.packets <- p[i]:
        default: //if pool full,throw away
        }
    }
}

func NewPacket() (p *RocPacket) {
    select {
    case p = <-pool.packets:
    default:
        p = newRocPacket()
    }
    return
}

func Payload(b []byte) *RocPacket {
    r := NewPacket()
    r.Write(b)
    return r
}

func PayloadIo(body io.ReadCloser) *RocPacket {
    r := NewPacket()
    _, _ = io.Copy(r.B, body)
    return r
}

func (r *RocPacket) recycle() {
    r.B.Reset()
}

func (r *RocPacket) Len() int {
    return r.B.Len()
}

func (r *RocPacket) Write(b []byte) {
    r.B.Write(b)
}

func (r *RocPacket) Bytes() []byte {
    return r.B.Bytes()
}

func (r *RocPacket) String() string {
    return r.B.String()
}
