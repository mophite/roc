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

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/internal/trace"
    "github.com/go-roc/roc/internal/trace/simple"
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

    //set http header
    Header map[string]string

    //setting http data context
    Data map[string]interface{}

    ////http writer
    //Writer http.ResponseWriter
    //
    ////http request
    //Request *http.Request
    //
    ////http request body
    //Body io.ReadCloser
}

func Background() *Context {
    return &Context{
        Trace:    simple.NewSimple(),
        Header:   make(map[string]string, 10),
        Data:     make(map[string]interface{}, 10),
        Metadata: new(metadata.Metadata),
    }
}

func (c *Context) SetHeader(key, value string) {
    c.Header[key] = value
}

func (c *Context) GetHeader(key string) string {
    return c.Header[key]
}

func (c *Context) Set(key string, value interface{}) {
    c.Data[key] = value
}

func (c *Context) Get(key string, value interface{}) {
    value = c.Data[key]
}

func (c *Context) WithMetadata(service, method string, meta map[string]string) error {
    m, err := metadata.EncodeMetadata(service, method, c.Trace.TraceId(), meta)
    if err != nil {
        return err
    }
    c.Metadata = m

    return nil
}

func NewContext(service, method, tracing string, meta map[string]string) (*Context, error) {
    m, err := metadata.EncodeMetadata(service, method, tracing, meta)
    if err != nil {
        return nil, err
    }
    return &Context{
        Metadata: m,
        Trace:    simple.WithSimple(tracing),
    }, nil
}

func FromMetadata(b []byte) *Context {
    m := metadata.DecodeMetadata(b)
    return &Context{
        Trace:       simple.WithSimple(m.Tracing()),
        Metadata:    m,
        ContentType: m.Get(namespace.DefaultHeaderContentType),
    }
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
