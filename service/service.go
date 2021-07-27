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

package service

import (
    "errors"
    "net/http"
    "os"
    "os/signal"
    "strings"

    "github.com/go-roc/roc/config"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service/opt"
    "github.com/gorilla/mux"

    "github.com/go-roc/roc/internal/etcd"
    "github.com/go-roc/roc/rlog/log"
    "github.com/go-roc/roc/service/client"
    "github.com/go-roc/roc/service/server"
)

type Service struct {
    //service options setting
    opts opt.Option

    //service exit channel
    exit chan struct{}

    //roc service client,for rpc to server
    client *client.Client

    //roc service server,listen and wait call
    server *server.Server
}

func New(opts ...opt.Options) *Service {
    s := &Service{
        opts: opt.NewOpts(opts...),
        exit: make(chan struct{}),
    }

    s.server = server.NewServer(s.opts)
    return s
}

func (s *Service) Client() *client.Client {
    if s.client == nil {
        s.client = client.NewClient(s.opts)
    }

    return s.client
}

func (s *Service) Server() *server.Server {
    return s.server
}

func (s *Service) Run() error {
    defer func() {
        if r := recover(); r != nil {
            rlog.Stack(r)
        }
    }()

    // handler signal
    ch := make(chan os.Signal)
    signal.Notify(ch, s.opts.Signal...)

    go func() {
        select {
        case c := <-ch:

            rlog.Infof("received signal %s [%s] transportServer exit!", c.String(), s.opts.Name)

            s.Close()

            for _, f := range s.opts.Exit {
                f()
            }

            s.exit <- struct{}{}
        }
    }()

    //run server
    s.server.Run()

    select {
    case <-s.exit:
    }

    return errors.New(s.opts.Name + " service is exit!")
}

func (s *Service) Close() {
    //close registry service discover
    if s.opts.Registry != nil {
        _ = s.opts.Registry.Deregister(s.opts.Endpoint)
        s.opts.Registry.Close()
        s.opts.Registry = nil
    }

    //close service client
    if s.client != nil {
        s.client.Close()
    }

    //close service server
    if s.server != nil {
        s.server.Close()
    }

    //close config setting
    config.Close()

    //close etcd client
    etcd.DefaultEtcd.Close()

    //todo flush rlog content
    log.Close()
}

var defaultRouter *mux.Router

// GROUP for not post http method
func (s *Service) GROUP(prefix string) *Service {
    if !strings.HasPrefix(prefix, "/") {
        prefix = "/" + prefix
    }
    if !strings.HasSuffix(prefix, "/") {
        prefix = prefix + "/"
    }

    //cannot had post prefix
    if strings.HasPrefix(prefix, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }

    defaultRouter = s.opts.Router.PathPrefix(prefix).Subrouter()
    return s
}

func (s *Service) GET(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodOptions, http.MethodGet)
}

func (s *Service) POST(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPost)
}

func (s *Service) PUT(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPut)
}

func (s *Service) DELETE(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodDelete)
}

func (s *Service) ANY(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler)
}

func (s *Service) HEAD(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodHead)
}

func (s *Service) PATCH(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPatch)
}

func (s *Service) CONNECT(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodConnect)
}

func (s *Service) TRACE(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodTrace)
}

func tidyRelativePath(relativePath string) string {
    //trim suffix "/"
    if strings.HasSuffix(relativePath, "/") {
        relativePath = strings.TrimSuffix(relativePath, "/")
    }

    //add prefix "/"
    if !strings.HasPrefix(relativePath, "/") {
        relativePath = "/" + relativePath
    }

    return relativePath
}

func GetApiPrefix() string {
    return DefaultApiPrefix
}
