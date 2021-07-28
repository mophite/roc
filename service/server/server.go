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
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/codec"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/opt"
    "github.com/go-roc/roc/service/router"
)

type Server struct {
    //run transportServer option
    opts opt.Option

    //transportServer exit channel
    exit chan struct{}

    //rpc transportServer router collection
    route *router.Router

    //api http server
    httpServer *http.Server
}

func NewServer(opts opt.Option) *Server {
    s := &Server{
        opts: opts,
        exit: make(chan struct{}),
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

            http.Handle(s.opts.ApiPrefix, s)

            s.httpServer = &http.Server{
                Handler:      s,
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

//roc not support method GET,because you can use other http web framework
//to build a restful api
//roc support POST,DELETE for compatible rrRouter ,witch request response way
//because ServeHTTP api need support json or proto data protocol
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    var c = context.Background()
    c.SetMethod(r.URL.Path)

    for k, v := range r.Header {
        if len(v) == 0 {
            continue
        }
        c.SetHeader(k, v[0])
    }
    c.ContentType = c.GetHeader(namespace.DefaultHeaderContentType)

    for i := range s.opts.HttpAddress {
        err := s.opts.HttpMiddleware[i](w, r)
        if err != nil {
            b := s.opts.Err.Encode(codec.GetCodec(c.ContentType), 400, err)
            w.Write(b)
            return
        }
    }

    switch r.Method {
    case http.MethodPost, http.MethodDelete:

        if _, ok := codec.DefaultCodecs[c.ContentType]; !ok {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`400 BAD REQUEST`))
            return
        }

        var req, rsp = parcel.PayloadIo(r.Body), parcel.NewPacket()
        defer func() {
            parcel.Recycle(req, rsp)
        }()

        _ = r.Body.Close()

        err := s.route.RRProcess(c, req, rsp)

        if err == router.ErrNotFoundHandler {
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`404 NOT FOUND`))
            return
        }

        if len(rsp.Bytes()) > 0 {
            w.WriteHeader(http.StatusOK)
            w.Write(rsp.Bytes())
        }
    }

    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write([]byte(`METHOD NOT ALLOWED`))
    return
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
