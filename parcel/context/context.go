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

package context

import (
    "fmt"
    "net"
    "strings"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/internal/trace"
    "github.com/go-roc/roc/internal/trace/simple"
    "github.com/go-roc/roc/parcel/codec"
    "github.com/go-roc/roc/parcel/metadata"
    "github.com/go-roc/roc/rlog/log"
)

type Context struct {

    //rpc metadata
    *metadata.Metadata

    //Trace exists throughout the life cycle of the context
    //trace is request flow trace
    //it's will be from web client,or generated on initialize
    Trace trace.Trace

    //Content-Type
    ContentType string

    ////http writer
    //Writer http.ResponseWriter
    //
    ////http request
    //Request *http.Request
    //
    ////http request body
    //Body io.ReadCloser
    data map[string]interface{}

    IsPutFile bool

    RemoteAddr string
}

func Background() *Context {
    return &Context{
        Trace:    simple.NewSimple(),
        Metadata: metadata.MallocMetadata(),
        data:     make(map[string]interface{}, 10),
    }
}

func (c *Context) Codec() codec.Codec {
    return codec.CodecType(c.ContentType)
}

func (c *Context) Clone(service, method string, meta map[string]string) (*Context, error) {
    m, err := metadata.EncodeMetadata(service, method, c.Trace.TraceId(), meta)
    if err != nil {
        return nil, err
    }
    s := *c
    s.Metadata = m

    return &s, nil
}

func (c *Context) SetSetupData(value []byte) {
    c.data[namespace.DefaultHeaderSetup] = value
}

func (c *Context) GetSetupData() []byte {
    b, _ := c.data[namespace.DefaultHeaderSetup].([]byte)
    return b
}

func FromMetadata(b []byte) *Context {
    m := metadata.DecodeMetadata(b)
    c := &Context{
        Trace:    simple.WithTrace(m.Tracing()),
        Metadata: m,
    }
    c.ContentType = c.GetHeader(namespace.DefaultHeaderContentType)
    c.Trace.SpreadOnce()
    return c
}

func (c *Context) ClientIP() string {
    clientIP := c.GetHeader("X-Forwarded-For")
    if clientIP != "" {
        s := strings.Split(clientIP, ",")
        if len(s) > 0 {
            clientIP = strings.TrimSpace(s[0])
        }
        if clientIP == "" {
            clientIP = strings.TrimSpace(c.GetHeader("X-Real-Ip"))
        }

        if clientIP != "" {
            return clientIP
        }
    }

    if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.RemoteAddr)); err == nil {
        return ip
    }

    return ""
}

func (c *Context) Get(key string) interface{} {
    return c.data[key]
}

func (c *Context) Set(key string, value interface{}) {
    c.data[key] = value
}

func (c *Context) GetHeader(key string) string {
    return c.GetMeta(key)
}

func (c *Context) SetHeader(key, value string) {
    c.SetMeta(key, value)
}

func (c *Context) Debug(msg ...interface{}) {
    c.Trace.Carrier()
    log.Debug(c.Trace.TraceId() + " |" + fmt.Sprintln(msg...))
}

func (c *Context) Info(msg ...interface{}) {
    c.Trace.Carrier()
    log.Info(c.Trace.TraceId() + " |" + fmt.Sprintln(msg...))
}

func (c *Context) Error(msg ...interface{}) {
    c.Trace.Carrier()
    log.Error(c.Trace.TraceId() + " |" + fmt.Sprintln(msg...))
}

func (c *Context) Debugf(f string, msg ...interface{}) {
    c.Trace.Carrier()
    log.Debug(c.Trace.TraceId() + " |" + fmt.Sprintf(f+"\n", msg...))
}

func (c *Context) Infof(f string, msg ...interface{}) {
    c.Trace.Carrier()
    log.Info(c.Trace.TraceId() + " |" + fmt.Sprintf(f+"\n", msg...))
}

func (c *Context) Errorf(f string, msg ...interface{}) {
    c.Trace.Carrier()
    log.Error(c.Trace.TraceId() + " |" + fmt.Sprintf(f+"\n", msg...))
}
