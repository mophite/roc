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

package hello

import (
	"net/http"
	"time"

	"tutorials/internal/ipc"
	"tutorials/proto/phello"

	"github.com/go-roc/roc"
	"github.com/go-roc/roc/rlog"
)

type HelloWorld struct {
	opt    []roc.InvokeOptions
	client phello.HelloWorldClient
}

// NewHello new Hello and initialize it for rpc client
// opt is configurable when request.
func NewHello(s *roc.Service) *HelloWorld {
	return &HelloWorld{
		client: phello.NewHelloWorldClient(s),
		opt: []roc.InvokeOptions{
			roc.WithName("srv.hello"),
		},
	}
}

func (h *HelloWorld) SayHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	now := time.Now()
	rsp, err := ipc.SayHello(context.Background(), &phello.SayReq{Inc: 1})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	rlog.Infof("FROM hello server: %v |latency=%v ms ", rsp.Inc, time.Since(now).Milliseconds())

	w.Write([]byte("succuess"))
}

func (h *HelloWorld) Say(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) error {
	return nil
}
