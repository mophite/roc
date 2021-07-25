package server

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/parcel/metadata"
	"github.com/go-roc/roc/rlog"
	"github.com/go-roc/roc/service/handler"
	"github.com/go-roc/roc/service/router"
)

type Server struct {
	//run server option
	opts Option

	//server exit channel
	exit chan struct{}

	//rpc server router collection
	route *router.Router

	//api router
	*mux.Router
}

func NewServer(opts ...Options) *Server {
	s := &Server{
		opts:   newOpts(opts...),
		exit:   make(chan struct{}),
		Router: mux.NewRouter(),
	}

	s.route = router.NewRouter(s.opts.wrappers, s.opts.err)

	s.opts.server.Accept(s.route)

	return s
}

func (s *Server) GetApiPrefix() string {
	return s.opts.apiPrefix
}

func (s *Server) Codec() codec.Codec {
	return codec.DefaultCodec
}

func (s *Server) Run() error {
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

	//run http server
	if s.opts.httpAddress != "" {
		go func() {

			s.PathPrefix(s.opts.apiPrefix).Handler(s)

			srv := &http.Server{
				Handler:      s.Router,
				Addr:         s.opts.httpAddress,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
				IdleTimeout:  time.Second * 60,
			}

			if err := srv.ListenAndServe(); err != nil {
				rlog.Error(err)
				s.exit <- struct{}{}
			}
		}()
	}

	rlog.Infof(
		"[TCP:%s][WS:%s][HTTP:%s] start success!",
		endpoint.DefaultLocalEndpoint.Absolute,
		s.opts.wssAddress,
		s.opts.httpAddress,
	)

	err := s.register()
	if err != nil {
		panic(err)
	}

	select {
	case <-s.exit:
	}

	return errors.New(s.opts.name + " server is exit!")
}

func (s *Server) register() error {
	return s.opts.registry.Register(endpoint.DefaultLocalEndpoint)
}

func (s *Server) RegisterHandler(method string, rr handler.Handler) {
	s.route.RegisterHandler(method, rr)
}

func (s *Server) RegisterStreamHandler(method string, rs handler.StreamHandler) {
	s.route.RegisterStreamHandler(method, rs)
}

func (s *Server) RegisterChannelHandler(method string, rs handler.ChannelHandler) {
	s.route.RegisterChannelHandler(method, rs)
}

func (s *Server) Close() {
	if s.opts.server != nil {
		s.opts.server.Close()
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var c = context.Background()

		c.Metadata = new(metadata.Metadata)
		c.SetMethod(r.URL.Path)

		var req, rsp = parcel.Payloader(r.Body), parcel.NewPacket()
		defer func() {
			parcel.Recycle(req, rsp)
		}()

		_ = r.Body.Close()

		err := s.route.RRProcess(c, req, rsp)

		//only allow post request less version 1.1.x
		//because ServeHTTP api need support json or proto data protocol
		if err == router.ErrNotFoundHandler {
			//todo 404 service code config to service
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write(rsp.Bytes())
	} else {
		//todo 404 service code config to service
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
