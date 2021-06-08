package router

import (
	"errors"
	"sync"

	"github.com/gogo/protobuf/proto"

	"roc/parcel"
	"roc/parcel/codec"
	"roc/parcel/context"
	"roc/rlog"
)

var (
	errNotFoundHandler = errors.New("not found rrRoute")
)

type Router struct {
	sync.Mutex
	rrRoute     map[string]parcel.Handler
	rsRoute     map[string]parcel.StreamHandler
	rcRoute     map[string]parcel.ChannelHandler
	wrappers    []parcel.Wrapper
	errorPacket parcel.ErrorPackager
	cc          codec.Codec
}

func NewRouter(wrappers []parcel.Wrapper, err parcel.ErrorPackager, c codec.Codec) *Router {
	return &Router{
		rrRoute:     make(map[string]parcel.Handler),
		rsRoute:     make(map[string]parcel.StreamHandler),
		rcRoute:     make(map[string]parcel.ChannelHandler),
		wrappers:    wrappers,
		errorPacket: err,
		cc:          c,
	}
}

func (r *Router) Codec() codec.Codec {
	return r.cc
}

func (r *Router) Error() parcel.ErrorPackager {
	return r.errorPacket
}

func (r *Router) RegisterHandler(method string, rr parcel.Handler) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rrRoute[method]; ok {
		panic("this rrRoute is already exist:" + method)
	}
	r.rrRoute[method] = rr
}

func (r *Router) RegisterStreamHandler(method string, rs parcel.StreamHandler) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rsRoute[method]; ok {
		panic("this rsRoute is already exist:" + method)
	}

	r.rsRoute[method] = rs
}

func (r *Router) RegisterChannelHandler(service string, rc parcel.ChannelHandler) {
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
		return errNotFoundHandler
	}

	resp, err := rr(c, req, r.interrupt())
	if err != nil {
		return err
	}

	b, err := r.cc.Encode(resp)
	if err != nil {
		return err
	}

	rsp.Write(b)

	return nil
}

func (r *Router) RSProcess(c *context.Context, req *parcel.RocPacket) (chan proto.Message, chan error) {

	rs, ok := r.rsRoute[c.Method()]
	if !ok {
		return nil, nil
	}

	return rs(c, req)
}

func (r *Router) RCProcess(c *context.Context, req chan *parcel.RocPacket, errs chan error) (chan proto.Message, chan error) {

	rc, ok := r.rcRoute[c.Method()]
	if !ok {
		return nil, nil
	}

	return rc(c, req, errs)
}

func (r *Router) interrupt() parcel.Interceptor {
	return func(c *context.Context, req proto.Message, fire parcel.Fire) (proto.Message, error) {
		for i := range r.wrappers {
			err := r.wrappers[i](c)
			if err != nil {
				c.Errorf("wrappers err=%v", err)
				return nil, err
			}
		}

		rsp, err := fire(c, req)
		if err != nil {
			c.Errorf("fire err=%v |FROM=%s", err, req.String())
			return nil, err
		}
		c.Infof("FROM=%s |TO=%s", req.String(), rsp.String())
		return rsp, nil
	}
}

func (r *Router) List() {
	rlog.Info("registered router list:")
	for k := range r.rrRoute {
		rlog.Infof("---------------------------------- [%s]", k)
	}
}
