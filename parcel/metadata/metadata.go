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
    "errors"
    "fmt"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/x"
    "github.com/go-roc/roc/x/bytesbuffpool"
)

// RsocketRpcVersion rsocket-rpc version
const RsocketRpcVersion = uint16(1)

type Metadata struct {
    V  string            `json:"version"`
    S  string            `json:"service"`
    M  string            `json:"method"`
    T  string            `json:"trace"`
    A  string            `json:"address"`
    P  []byte            `json:"payload"`
    Md map[string]string `json:"meta"`
}

func MallocMetadata() *Metadata {
    return &Metadata{Md: make(map[string]string, 10)}
}

func EncodeMetadata(service, method, tracing string, meta map[string]string) (*Metadata, error) {
    b, err := encodeMetadata(service, method, tracing, meta)
    if err != nil {
        return nil, err
    }

    return DecodeMetadata(b)
}

func (m *Metadata) Payload() []byte {
    return m.P
}

func (m *Metadata) GetMeta(key string) string {
    return m.Md[key]
}

func (m *Metadata) SetMeta(key, value string) {
    m.Md[key] = value
}

func (m *Metadata) Service() string {
    return m.S
}

func (m *Metadata) Method() string {
    return m.M
}

func (m *Metadata) SetMethod(method string) {
    m.M = method
}

func (m *Metadata) Version() string {
    return m.V
}

func (m *Metadata) Tracing() string {
    return m.T
}

func (m *Metadata) Address() string {
    return m.A
}

func (m *Metadata) String() string {
    var tr string
    if b := m.Tracing(); len(b) < 1 {
        tr = "<nil>"
    } else {
        tr = "0x" + hex.EncodeToString([]byte(b))
    }

    var s string
    if b := m.Md; len(b) < 1 {
        s = "<nil>"
    } else {
        s = "0x" + hex.EncodeToString(m.getMetadata())
    }
    return fmt.Sprintf(
        "Metadata{version=%s, service=%s, method=%s, tracing=%s, metadata=%s}",
        m.Version(),
        m.Service(),
        m.Method(),
        tr,
        s,
    )
}

func (m *Metadata) VersionUint16() uint16 {
    return binary.BigEndian.Uint16(m.P)
}

func (m *Metadata) getService() string {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2

    return string(m.P[offset : offset+serviceLen])
}

func (m *Metadata) getMethod() string {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2 + serviceLen

    methodLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2

    return string(m.P[offset : offset+methodLen])
}

func (m *Metadata) getTrace() []byte {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2 + serviceLen

    methodLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2 + methodLen

    tracingLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2

    if tracingLen > 0 {
        return m.P[offset : offset+tracingLen]
    } else {
        return nil
    }
}

func (m *Metadata) getMetadata() []byte {
    offset := 2

    serviceLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2 + serviceLen

    methodLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2 + methodLen

    tracingLen := int(binary.BigEndian.Uint16(m.P[offset : offset+2]))
    offset += 2 + tracingLen

    return m.P[offset:]
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

func DecodeMetadata(payload []byte) (*Metadata, error) {

    m := &Metadata{P: payload}

    err := x.Jsoniter.Unmarshal(m.getMetadata(), &m.Md)
    if err != nil {
        return nil, err
    }

    m.M = m.getMethod()
    if m.M == "" {
        return nil, errors.New("no method")
    }

    m.S = m.getService()

    m.T = x.BytesToString(m.getTrace())

    m.V = m.GetMeta(namespace.DefaultHeaderVersion)
    if m.V == "" {
        m.V = namespace.DefaultVersion
    }

    m.A = m.GetMeta(namespace.DefaultHeaderAddress)

    return m, nil
}
