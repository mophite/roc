package client

import (
	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/rlog/log"
	"github.com/go-roc/roc/service/strategy"
	"github.com/gogo/protobuf/proto"
)

type Client struct {
	//run server option
	opts Option

	//server exit channel
	exit chan struct{}

	//The strategy used by the client to initiate an rpc request, with roundRobin or direct ip request
	strategy strategy.Strategy
}

func NewService(opts ...Options) *Client {
	s := &Client{
		opts: newOpts(opts...),
		exit: make(chan struct{}),
	}

	s.strategy = strategy.NewStrategy(s.opts.e, s.opts.registry, s)

	return s
}

func (s *Client) Codec() codec.Codec {
	return s.opts.cc
}

// InvokeRR rpc request requestResponse,it's block request,one request one response
func (s *Client) InvokeRR(c *Context, method string, req, rsp proto.Message, opts ...InvokeOptions) error {

	// new a invoke setting
	invoke, err := newInvoke(c, method, s, opts...)
	if err != nil {
		return err
	}

	var conn *Conn

	// if address is nil ,user roundRobin strategy
	// otherwise straight to invoke ip server
	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}

	if err != nil {
		return err
	}

	return invoke.invokeRR(c, req, rsp, conn, s.opts)
}

// InvokeRS rpc request requestStream,it's one request and multiple response
func (s *Client) InvokeRS(c *Context, method string, req proto.Message, opts ...InvokeOptions) (
	chan []byte,
	chan error,
) {

	// new a invoke setting
	invoke, err := newInvoke(c, method, s, opts...)
	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	var conn *Conn

	// if address is nil ,user roundRobin strategy
	// otherwise straight to invoke ip server
	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}

	//encode req body to roc packet
	b, err := s.opts.cc.Encode(req)

	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	return conn.Client().RS(c, parcel.Payload(b))
}

// InvokeRC rpc request requestChannel,it's multiple request and multiple response
func (s *Client) InvokeRC(
	c *Context,
	method string,
	req chan []byte,
	errIn chan error,
	opts ...InvokeOptions,
) (chan []byte, chan error) {

	// new a invoke setting
	invoke, err := newInvoke(c, method, s, opts...)
	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	var conn *Conn

	// if address is nil ,user roundRobin strategy
	// otherwise straight to invoke ip server
	if invoke.opts.address != "" {
		conn, err = invoke.strategy.Straight(invoke.opts.scope, invoke.opts.address)
	} else {
		conn, err = invoke.strategy.Next(invoke.opts.scope)
	}
	if err != nil {
		// create a chan error response
		var errs = make(chan error)
		errs <- err
		close(errs)
		return nil, errs
	}

	return conn.Client().RC(c, req, errIn)
}

func (s *Client) Close() {
	if s.opts.registry != nil {
		_ = s.opts.registry.Deregister(s.opts.e)
		s.opts.registry.Close()
	}

	if s.strategy != nil {
		s.strategy.Close()
	}

	//todo flush rlog content
	log.Close()
}
