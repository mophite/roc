package server

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/service/handler"
	"github.com/go-roc/roc/service/router"
	"github.com/go-roc/roc/x"
	"github.com/gorilla/mux"

	"github.com/go-roc/roc/rlog"
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

func NewService(opts ...Options) *Server {
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
		s.opts.e.Absolute,
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
	return registry.DefaultRegistry.Register(s.opts.e)
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

	h, ok := s.route.ApiProcess(r.URL.Path)

	//only allow post request less version 1.1.x
	//because ServeHTTP api need support json or proto data protocol
	if !ok {
		//todo 404 service code config to service
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		var c = context.Background()
		c.Body = r.Body

		rsp, err := h(c)

		_ = r.Body.Close()

		if err != nil {
			//todo 500 service code config to service
			rlog.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(x.MustMarshal(rsp))
	}
}

func (s *Server) RegisterApiRouter(relativePath string, apiHandler handler.ApiRocHandler) {
	s.route.RegisterApiHandler(relativePath, apiHandler)
}
