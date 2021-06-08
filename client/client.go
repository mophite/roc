package client

import (
	"roc/parcel"
	"roc/parcel/codec"
	"time"

	"github.com/gogo/protobuf/proto"

	"roc/internal/registry"
	"roc/internal/transport"
	"roc/parcel/context"
)

type RocClient struct {
	opts     option
	strategy Strategy
}

func (r *RocClient) Registry() registry.Registry {
	return r.opts.registry
}

func (r *RocClient) Client() transport.Client {
	return r.opts.client
}

func (r *RocClient) ConnectTimeout() time.Duration {
	return r.opts.connectTimeout
}

func (r *RocClient) KeepaliveInterval() time.Duration {
	return r.opts.keepaliveInterval
}

func (r *RocClient) KeepaliveLifetime() time.Duration {
	return r.opts.keepaliveLifetime
}

func (r *RocClient) Codec() codec.Codec {
	return r.opts.cc
}

func NewRocClient(opts ...Options) (lp *RocClient) {

	lp = &RocClient{opts: newOpts(opts...)}
	lp.strategy = newStrategy(lp.opts.registry, lp.opts.client)

	return
}

func (r *RocClient) InvokeRR(c *context.Context, method string, req, rsp proto.Message, opts ...InvokeOptions) error {

	invoke, err := newInvoke(c, method, r, opts...)
	if err != nil {
		return err
	}

	var conn *Conn

	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}
	if err != nil {
		return err
	}

	return invoke.invokeRR(c, req, rsp, conn, r.opts)
}

func (r *RocClient) InvokeRS(c *context.Context, method string, req proto.Message, opts ...InvokeOptions) (chan []byte, chan error) {

	invoke, err := newInvoke(c, method, r, opts...)
	if err != nil {
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	var conn *Conn

	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}

	b, err := r.opts.cc.Encode(req)

	if err != nil {
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	return conn.Client().RS(c, parcel.Payload(b))
}

func (r *RocClient) InvokeRC(c *context.Context, method string, req chan []byte, errIn chan error, opts ...InvokeOptions) (chan []byte, chan error) {

	invoke, err := newInvoke(c, method, r, opts...)
	if err != nil {
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	var conn *Conn

	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}
	if err != nil {
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	return conn.Client().RC(c, req, errIn)
}

func (r *RocClient) Close() {
	r.strategy.Close()
}
