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
	"sync"
)

var packetPool = sync.Pool{New: func() interface{} {
	return &RocPacket{B: new(bytes.Buffer)}
}}

type RocPacket struct {
	B *bytes.Buffer
}

func Recycle(p ...*RocPacket) {
	for i := range p {
		p[i].B.Reset()
		packetPool.Put(p[i])
	}
}

func NewPacket() *RocPacket {
	return packetPool.Get().(*RocPacket)
}

func Payload(b []byte) *RocPacket {
	r := packetPool.Get().(*RocPacket)
	r.Write(b)
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
