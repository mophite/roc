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
    "bytes"
    ctx "context"
    "io"
    "net/http"
    "strings"
    "sync"
    "time"

    "github.com/rs/cors"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/parcel/packet"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/service/opt"
    "github.com/go-roc/roc/service/router"
)

type Server struct {

    //wait for server init
    wg *sync.WaitGroup

    //run transportServer option
    opts opt.Option

    //transportServer exit channel
    exit chan struct{}

    //rpc transportServer router collection
    route *router.Router

    //api http server
    httpServer *http.Server
}

func (s *Server) Name() string {
    name := s.opts.Name
    ss := strings.Split(name, ".")

    if len(ss) > 1 {
        name = ss[len(ss)-1]
    }

    return name
}

func NewServer(opts opt.Option) *Server {
    s := &Server{
        wg:   new(sync.WaitGroup),
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

    s.opts.TransportServer.Run(s.wg)

    //run http transportServer
    if s.opts.HttpAddress != "" {
        go func() {

            prefix := s.opts.Name

            if !strings.HasPrefix(prefix, "/") {
                prefix = "/" + prefix
            }

            if !strings.HasSuffix(prefix, "/") {
                prefix = prefix + "/"
            }

            http.Handle(prefix, cors.New(*s.opts.CorsOptions).Handler(s))

            s.httpServer = &http.Server{
                Handler:      s,
                Addr:         s.opts.HttpAddress,
                WriteTimeout: 15 * time.Second,
                ReadTimeout:  15 * time.Second,
                IdleTimeout:  time.Second * 60,
            }

            if err := s.httpServer.ListenAndServe(); err != nil {
                rlog.Errorf("service %s |err=%v", s.opts.Name, err)
            }
        }()
    }

    s.wg.Wait()

    rlog.Infof(
        "[TCP:%s:%d][WS:%s][HTTP:%s] start success!",
        s.opts.LocalIp, s.opts.TcpAddress,
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

//roc don't suggest method like GET,because you can use other http web framework
//to build a restful api with not by roc
//roc support POST,DELETE,PUT,GET,OPTIONS for compatible rrRouter ,witch request response way
//because ServeHTTP api need support json or proto data protocol
//suggest just use POST,PUT for your roc service
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    var c = context.Background()
    c.SetMethod(r.URL.Path)

    for k, v := range r.Header {
        if len(v) == 0 {
            continue
        }
        c.SetHeader(k, v[0])
    }
    c.ContentType = c.GetHeader(namespace.DefaultHeaderContentType)
    c.SetCodec()

    c.RemoteAddr = r.RemoteAddr

    w.Header().Set(namespace.DefaultHeaderContentType, c.ContentType)
    w.Header().Set(namespace.DefaultHeaderTrace, c.Trace.TraceId())

    for i := range s.opts.Dog {
        rsp, err := s.opts.Dog[i](c)
        if err != nil {
            c.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(c.Codec().MustEncode(rsp))
            return
        }
    }

    switch r.Method {
    case http.MethodPost, http.MethodDelete:

        var req, rsp = parcel.PayloadIo(r.Body), parcel.NewPacket()
        defer func() {
            parcel.Recycle(req)
            parcel.Recycle(rsp)
        }()

        _ = r.Body.Close()

        err := s.route.RRProcess(c, req, rsp)

        if err == router.ErrNotFoundHandler {
            c.Errorf("path=%s |service=%s", r.URL.Path, s.opts.Name)
            w.WriteHeader(http.StatusNotFound)
            w.Write(s.opts.Err.Error404(c))
            return
        }

        if len(rsp.Bytes()) > 0 {
            w.WriteHeader(http.StatusOK)
            w.Write(rsp.Bytes())
        } else if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusInternalServerError)
            w.Write(s.opts.Err.Error500(c))
        }

        return

    case http.MethodPut:

        f, h, err := r.FormFile("file")
        if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(s.opts.Err.Error400(c))
            return
        }

        var buf = bytes.NewBuffer(make([]byte, 0, 10485760))

        io.Copy(buf, f)

        var fileReq = &packet.FileReq{}
        fileReq.Body = buf.Bytes()
        fileReq.FileSize = h.Size
        fileReq.FileName = h.Filename
        fileReq.Extra = r.FormValue("extra")

        fb, err := c.Codec().Encode(fileReq)
        if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(s.opts.Err.Error400(c))
            return
        }

        var req, rsp = parcel.Payload(fb), parcel.NewPacket()
        defer func() {
            parcel.Recycle(req)
            parcel.Recycle(rsp)
        }()
        _ = r.Body.Close()

        c.IsPutFile = true
        err = s.route.RRProcess(c, req, rsp)

        if err == router.ErrNotFoundHandler {
            c.Errorf("err=%v |path=%s", err, c.Method())
            w.Header().Set("Content-type", "text/plain")
            w.WriteHeader(http.StatusNotFound)
            w.Write(s.opts.Err.Error404(c))
            return
        }

        if len(rsp.Bytes()) > 0 {
            w.WriteHeader(http.StatusOK)
            w.Write(rsp.Bytes())
        } else if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusInternalServerError)
            w.Write(s.opts.Err.Error500(c))
        }

        return

    case http.MethodGet:

        values := r.URL.Query()

        var apiReq = &packet.ApiReq{Params: make(map[string]string, len(values))}
        for k, v := range values {
            if len(v) > 0 {
                apiReq.Params[k] = v[0]
            }
        }

        fb, err := c.Codec().Encode(apiReq)
        if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(s.opts.Err.Error400(c))
            return
        }

        var req, rsp = parcel.Payload(fb), parcel.NewPacket()
        defer func() {
            parcel.Recycle(req)
            parcel.Recycle(rsp)
        }()

        err = s.route.RRProcess(c, req, rsp)

        if err == router.ErrNotFoundHandler {
            c.Errorf("err=%v |path=%s", err, c.Method())
            w.WriteHeader(http.StatusNotFound)
            w.Write(s.opts.Err.Error404(c))
            return
        }

        if len(rsp.Bytes()) > 0 {
            w.WriteHeader(http.StatusOK)
            w.Write(rsp.Bytes())
        } else if err != nil {
            rlog.Error(err)
            w.Header().Set("Content-type", "text/plain")
            w.WriteHeader(http.StatusInternalServerError)
            w.Write(s.opts.Err.Error500(c))
        }

        return

    case http.MethodOptions:
        w.WriteHeader(http.StatusOK)
        return
    }

    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write(s.opts.Err.Error405(c))
    return
}

func (s *Server) CloseServer() {
    cc, cancel := ctx.WithTimeout(ctx.Background(), time.Second*5)
    defer cancel()

    if s.httpServer != nil {
        _ = s.httpServer.Shutdown(cc)
    }

    if s.opts.TransportServer != nil {
        s.opts.TransportServer.Close()
    }
}
