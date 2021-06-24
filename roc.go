package roc

import (
	"errors"
	"os"
	"os/signal"

	"github.com/gogo/protobuf/proto"

	"github.com/go-roc/roc/internal/router"
	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/rlog"
	"github.com/go-roc/roc/rlog/log"
)

type Service struct {
	//run server option
	opts Option

	//server exit channel
	exit chan struct{}

	//server router collection
	route *router.Router

	//The strategy used by the client to initiate an rpc request, with roundRobin or direct ip request
	strategy Strategy
}

func NewService(opts ...Options) *Service {
	s := &Service{
		opts: newOpts(opts...),
		exit: make(chan struct{}),
	}

	s.route = router.NewRouter(s.opts.wrappers, s.opts.err, s.opts.cc)

	//NOTICE: don't register wss to sd.
	//if r.SetupWss() {
	//	err := r.opts.register.Register(r.WssAddress(), "wss")
	//	if err != nil {
	//		return nil
	//	}
	//}
	s.opts.server.Accept(s.route)

	s.strategy = newStrategy(s.opts.registry, s.opts.client)
	return s
}

func (s *Service) Codec() codec.Codec {
	return s.opts.cc
}

// InvokeRR rpc request requestResponse,it's block request,one request one response
func (s *Service) InvokeRR(c *context.Context, method string, req, rsp proto.Message, opts ...InvokeOptions) error {

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
func (s *Service) InvokeRS(c *context.Context, method string, req proto.Message, opts ...InvokeOptions) (chan []byte, chan error) {

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
func (s *Service) InvokeRC(c *context.Context, method string, req chan []byte, errIn chan error, opts ...InvokeOptions) (chan []byte, chan error) {

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

func (s *Service) Run() error {
	defer func() {
		if r := recover(); r != nil {
			rlog.Stack(r)
		}
	}()

	// handler signal
	ch := make(chan os.Signal)
	signal.Notify(ch, s.opts.signal...)

	go func() {
		select {
		case c := <-ch:

			rlog.Infof("received signal %s [%s] server exit!", c.String(), s.opts.name)

			s.Close()

			for _, f := range s.opts.exit {
				f()
			}

			s.exit <- struct{}{}
		}
	}()

	// echo method list
	s.route.List()
	s.opts.server.Run()

	rlog.Infof("[tcp:%s] AND [ws:%s] is start success!",
		s.opts.e.Absolute,
		s.opts.wssAddress,
	)

	err := s.register()
	if err != nil {
		panic(err)
	}

	select { case <-s.exit: }

	return errors.New(s.opts.name + " server is exit!")
}

func (s *Service) register() error {
	return s.opts.registry.Register(s.opts.e)
}

func (s *Service) RegisterHandler(method string, rr parcel.Handler) {
	s.route.RegisterHandler(method, rr)
}

func (s *Service) RegisterStreamHandler(method string, rs parcel.StreamHandler) {
	s.route.RegisterStreamHandler(method, rs)
}

func (s *Service) RegisterChannelHandler(method string, rs parcel.ChannelHandler) {
	s.route.RegisterChannelHandler(method, rs)
}

func (s *Service) Close() {
	if s.opts.registry != nil {
		_ = s.opts.registry.Deregister(s.opts.e)
		s.opts.registry.Close()
	}

	if s.strategy != nil {
		s.strategy.Close()
	}

	if s.opts.server != nil {
		s.opts.server.Close()
	}

	if s.opts.client != nil {
		s.opts.client.Close()
	}

	//todo flush rlog content
	log.Close()
}
