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

package server

import (
    ctx "context"
    "net/http"
    "time"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/service/opt"
    "github.com/gorilla/mux"

    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/parcel/metadata"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/router"
)

type Server struct {
    //run transportServer option
    opts opt.Option

    //transportServer exit channel
    exit chan struct{}

    //rpc transportServer router collection
    route *router.Router

    //api router
    *mux.Router

    //api http server
    httpServer *http.Server
}

func NewServer(opts opt.Option) *Server {
    s := &Server{
        opts:   opts,
        exit:   make(chan struct{}),
        Router: mux.NewRouter(),
    }

    s.route = router.NewRouter(s.opts.Wrappers, s.opts.Err)

    s.opts.TransportServer.Accept(s.route)

    return s
}

func (s *Server) Run() {
    // echo method list
    s.route.List()

    s.opts.TransportServer.Run()

    //run http transportServer
    if s.opts.HttpAddress != "" {
        go func() {

            s.PathPrefix(s.opts.ApiPrefix).Handler(s)

            s.httpServer = &http.Server{
                Handler:      s.Router,
                Addr:         s.opts.HttpAddress,
                WriteTimeout: 15 * time.Second,
                ReadTimeout:  15 * time.Second,
                IdleTimeout:  time.Second * 60,
            }

            if err := s.httpServer.ListenAndServe(); err != nil {
                rlog.Error(err)
            }
        }()
    }

    rlog.Infof(
        "[TCP:%s][WS:%s][HTTP:%s] start success!",
        s.opts.TcpAddress,
        s.opts.WssAddress,
        s.opts.HttpAddress,
    )

    err := s.register()
    if err != nil {
        panic(err)
    }
}

func (s *Server) register() error {
    return s.opts.Registry.Register(s.opts.Endpoint)
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        var c = context.Background()

        c.Metadata = new(metadata.Metadata)
        c.SetMethod(r.URL.Path)
        c.ContentType = r.Header.Get(namespace.DefaultHeaderContentType)

        var req, rsp = parcel.PayloadIo(r.Body), parcel.NewPacket()
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

        if err != nil {
            c.Error(err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        w.Write(rsp.Bytes())
    } else {
        //todo service code config to service
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
}

func (s *Server) Close() {
    cc, cancel := ctx.WithTimeout(ctx.Background(), time.Second*5)
    defer cancel()

    if s.httpServer != nil {
        _ = s.httpServer.Shutdown(cc)
    }

    if s.opts.TransportServer != nil {
        s.opts.TransportServer.Close()
    }
}
