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

// Package metadata for websocket or socket
// from rsocket-rpc-go Metadata
package metadata

import (
    "encoding/binary"
    "encoding/hex"
    "fmt"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/x"
    "github.com/go-roc/roc/x/bytesbuffpool"
)

// RsocketRpcVersion rsocket-rpc version
const RsocketRpcVersion = uint16(1)

type Metadata struct {
    version string
    service string
    method  string
    trace   string
    address string
    meta    map[string]string
    payload []byte
}

func MallocMetadata() *Metadata {
    return &Metadata{meta: make(map[string]string, 10)}
}

func DecodeMetadata(b []byte) *Metadata {
    return decodeMetadata(b)
}

func EncodeMetadata(service, method, tracing string, meta map[string]string) (*Metadata, error) {
    b, err := encodeMetadata(service, method, tracing, meta)
    if err != nil {
        return nil, err
    }

    return decodeMetadata(b), nil
}

func (p *Metadata) Payload() []byte {
    return p.payload
}

func (p *Metadata) Set(key, value string) {
    p.meta[key] = value
}

func (p *Metadata) Get(key string) string {
    return p.meta[key]
}

func (p *Metadata) Service() string {
    return p.service
}

func (p *Metadata) Method() string {
    return p.method
}

func (p *Metadata) SetMethod(method string) {
    p.method = method
}

func (p *Metadata) Version() string {
    return p.version
}

func (p *Metadata) Tracing() string {
    return p.trace
}

func (p *Metadata) Address() string {
    return p.address
}

func (p *Metadata) String() string {
    var tr string
    if b := p.Tracing(); len(b) < 1 {
        tr = "<nil>"
    } else {
        tr = "0x" + hex.EncodeToString([]byte(b))
    }

    var m string
    if b := p.meta; len(b) < 1 {
        m = "<nil>"
    } else {
        m = "0x" + hex.EncodeToString(p.getMetadata())
    }
    return fmt.Sprintf(
        "Metadata{version=%s, service=%s, method=%s, tracing=%s, metadata=%s}",
        p.Version(),
        p.Service(),
        p.Method(),
        tr,
        m,
    )
}

func (p *Metadata) VersionUint16() uint16 {
    return binary.BigEndian.Uint16(p.payload)
}

func (p *Metadata) getService() string {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2

    return string(p.payload[offset : offset+serviceLen])
}

func (p *Metadata) getMethod() string {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2 + serviceLen

    methodLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2

    return string(p.payload[offset : offset+methodLen])
}

func (p *Metadata) getTrace() []byte {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2 + serviceLen

    methodLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2 + methodLen

    tracingLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2

    if tracingLen > 0 {
        return p.payload[offset : offset+tracingLen]
    } else {
        return nil
    }
}

func (p *Metadata) getMetadata() []byte {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2 + serviceLen

    methodLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2 + methodLen

    tracingLen := int(binary.BigEndian.Uint16(p.payload[offset : offset+2]))
    offset += 2 + tracingLen

    return p.payload[offset:]
}

func encodeMetadata(service, method, tracing string, metadata map[string]string) (m []byte, err error) {

    w := bytesbuffpool.Get()
    // write version
    err = binary.Write(w, binary.BigEndian, RsocketRpcVersion)
    if err != nil {
        return
    }
    // write service
    err = binary.Write(w, binary.BigEndian, uint16(len(service)))
    if err != nil {
        return
    }
    _, err = w.WriteString(service)
    if err != nil {
        return
    }
    // write method
    err = binary.Write(w, binary.BigEndian, uint16(len(method)))
    if err != nil {
        return
    }
    _, err = w.WriteString(method)
    if err != nil {
        return
    }
    // write tracing
    lenTracing := uint16(len(tracing))
    err = binary.Write(w, binary.BigEndian, lenTracing)
    if err != nil {
        return
    }
    if lenTracing > 0 {
        _, err = w.WriteString(tracing)
        if err != nil {
            return
        }
    }
    // write metadata
    if l := len(metadata); l > 0 {
        _, err = w.Write(x.MustMarshal(metadata))
        if err != nil {
            return
        }
    }
    m = w.Bytes()

    bytesbuffpool.Put(w)
    return
}

func decodeMetadata(payload []byte) *Metadata {

    m := &Metadata{payload: payload}

    x.MustUnmarshal(m.getMetadata(), &m.meta)

    m.method = m.getMethod()
    m.service = m.getService()
    m.trace = x.BytesToString(m.getTrace())
    m.version = m.Get(namespace.DefaultHeaderVersion)
    m.address = m.Get(namespace.DefaultHeaderAddress)

    return m
}
