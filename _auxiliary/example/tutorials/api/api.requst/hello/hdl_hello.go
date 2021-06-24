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

	"github.com/coreos/etcd/clientv3"

	"github.com/go-roc/roc"
	"github.com/go-roc/roc/_auxiliary/example/tutorials/proto/pbhello"
	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/rlog"
)

type Hello struct {
	opt    roc.InvokeOptions
	client pbhello.HelloWorldClient
}

// NewHello new Hello and initialize it for rpc client
// opt is configurable when request.
func NewHello() *Hello {
	return &Hello{
		client: pbhello.NewHelloWorldClient(
			roc.NewService(
				roc.TCPAddress("127.0.0.1:8899"),
				roc.Namespace("srv.hello"),
				roc.EtcdConfig(&clientv3.Config{
					Endpoints: []string{"82.157.14.79:2379"},
				}),
			),
		),
		opt: roc.WithName("srv.hello"),
	}
}

func (h *Hello) SayHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	rsp, err := h.client.Say(context.Background(), &pbhello.SayReq{Inc: 1}, h.opt)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	rlog.Info("FROM helloe server: ", rsp.Inc)

	w.Write([]byte("succuess"))
}
