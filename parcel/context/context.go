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
	"io"
	"net/http"

	"github.com/go-roc/roc/internal/trace"
	"github.com/go-roc/roc/internal/trace/simple"
	"github.com/go-roc/roc/parcel/metadata"
	"github.com/go-roc/roc/rlog/log"
)

type Context struct {
	*metadata.Metadata

	//tracker exists throughout the life cycle of the context
	Trace trace.Trace

	//http writer
	Writer http.ResponseWriter

	//http request
	Request *http.Request

	//http request body
	Body io.ReadCloser
}

func Background() *Context {
	return new(Context)
}

func (c *Context) WithMetadata(service, method, tracing string, meta map[string]string) error {
	m, err := metadata.EncodeMetadata(service, method, tracing, meta)
	if err != nil {
		return err
	}
	c.Metadata = m

	if tracing == "" {
		c.Trace = simple.NewSimple(tracing)
	}

	return nil
}

func NewContext(service, method, tracing string, meta map[string]string) (*Context, error) {
	m, err := metadata.EncodeMetadata(service, method, tracing, meta)
	if err != nil {
		return nil, err
	}
	return &Context{
		Metadata: m,
		Trace:    simple.NewSimple(tracing),
	}, nil
}

func FromMetadata(b []byte) *Context {
	m := metadata.DecodeMetadata(b)
	return &Context{
		Trace:    simple.NewSimple(m.Tracing()),
		Metadata: m,
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
