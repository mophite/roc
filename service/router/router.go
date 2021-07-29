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

package router

import (
	"errors"
	"sync"

	"github.com/gogo/protobuf/proto"

	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/service/handler"

	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/rlog"
)

var (
	ErrNotFoundHandler = errors.New("not found route path")
)

type Router struct {
	sync.Mutex
	//requestResponse map cache handler
	rrRoute map[string]handler.Handler

	//requestStream map cache streamHandler
	rsRoute map[string]handler.StreamHandler

	//requestChannel map cache channelHandler
	rcRoute map[string]handler.ChannelHandler

	//wrapper middleware
	wrappers []handler.WrapperHandler

	//configurable error message return
	errorPacket parcel.ErrorPackager
}

// NewRouter create a new Router
func NewRouter(wrappers []handler.WrapperHandler, err parcel.ErrorPackager) *Router {
	return &Router{
		rrRoute:     make(map[string]handler.Handler),
		rsRoute:     make(map[string]handler.StreamHandler),
		rcRoute:     make(map[string]handler.ChannelHandler),
		wrappers:    wrappers,
		errorPacket: err,
	}
}

func (r *Router) Error() parcel.ErrorPackager {
	return r.errorPacket
}

func (r *Router) RegisterHandler(method string, rr handler.Handler) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rrRoute[method]; ok {
		panic("this rrRoute is already exist:" + method)
	}
	r.rrRoute[method] = rr
}

func (r *Router) RegisterStreamHandler(method string, rs handler.StreamHandler) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rsRoute[method]; ok {
		panic("this rsRoute is already exist:" + method)
	}

	r.rsRoute[method] = rs
}

func (r *Router) RegisterChannelHandler(service string, rc handler.ChannelHandler) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rcRoute[service]; ok {
		panic("this rcRoute is already exist:" + service)
	}

	r.rcRoute[service] = rc
}

func (r *Router) RRProcess(c *context.Context, req *parcel.RocPacket, rsp *parcel.RocPacket) error {
	rr, ok := r.rrRoute[c.Method()]
	if !ok {
		return ErrNotFoundHandler
	}
	resp, err := rr(c, req, r.interrupt())
	if resp != nil {
		b, err := codec.GetCodec(c.ContentType).Encode(resp)
		if err != nil {
			c.Error(err)
			return err
		}

		rsp.Write(b)
	}

	if err != nil {
		c.Error(err)
		return err
	}

	return nil
}

func (r *Router) RSProcess(c *context.Context, req *parcel.RocPacket) (chan proto.Message, chan error) {

	// interrupt
	for i := range r.wrappers {
		_, err := r.wrappers[i](c)
		if err != nil {
			c.Errorf("wrappers err=%v", err)
			var errs = make(chan error)
			errs <- err
			close(errs)
			return nil, errs
		}
	}

	rs, ok := r.rsRoute[c.Method()]
	if !ok {
		return nil, nil
	}

	return rs(c, req)
}

func (r *Router) RCProcess(c *context.Context, req chan *parcel.RocPacket, errs chan error) (
	chan proto.Message,
	chan error,
) {
	// interrupt when occur error
	for i := range r.wrappers {
		_, err := r.wrappers[i](c)
		if err != nil {
			c.Errorf("wrappers err=%v", err)
			var errs = make(chan error)
			errs <- err
			close(errs)
			return nil, errs
		}
	}

	rc, ok := r.rcRoute[c.Method()]
	if !ok {
		return nil, nil
	}

	return rc(c, req, errs)
}

func (r *Router) interrupt() handler.Interceptor {
	return func(c *context.Context, req proto.Message, fire handler.Fire) (proto.Message, error) {
		// interrupt when occur error
		for i := range r.wrappers {
			rsp, err := r.wrappers[i](c)
			if err != nil {
				c.Errorf("wrappers err=%v", err)
				return rsp, err
			}
		}

		rsp, err := fire(c, req)
		if err != nil {
			c.Errorf("fire err=%v |FROM=%s", err, req.String())
			return rsp, err
		}
		c.Infof("FROM=%s |TO=%s", req.String(), rsp.String())
		return rsp, nil
	}
}

func (r *Router) List() {
	rlog.Info("registered Router list:")
	for k := range r.rrRoute {
		rlog.Infof("-------------------RR------------------- [%s]", k)
	}
}
