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
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
	"github.com/go-roc/roc/parcel/context"
	"github.com/go-roc/roc/rlog"
)

type Hello struct{}

func (h *Hello) SaySrv(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) {
	rsp.Pong = "pong"
}

func (h *Hello) SayStream(c *context.Context, req *phello.SayReq) chan *phello.SayRsp {
	var rsp = make(chan *phello.SayRsp)

	go func() {
		var count uint32
		for i := 0; i < 3; i++ {
			rsp <- &phello.SayRsp{Pong: strconv.Itoa(i)}
			atomic.AddUint32(&count, 1)
			time.Sleep(time.Second * 1)
		}

		c.Info("say stream example count is: ", atomic.LoadUint32(&count))

		close(rsp)
	}()

	return rsp
}

func (h *Hello) SayChannel(c *context.Context, req chan *phello.SayReq, exit chan struct{}) chan *phello.SayRsp {
	var rsp = make(chan *phello.SayRsp)

	go func() {
	QUIT:
		for {
			select {
			case data, ok := <-req:
				if !ok {
					break QUIT
				}

				rlog.Infof("FROM |req=%s", data.String())
				//test channel sending frequency
				//time.Sleep(time.Second)
				rsp <- &phello.SayRsp{Pong: "pong"}
			case <-exit:
				rlog.Info("exit")
				break QUIT
			}
		}

		close(rsp)
	}()

	return rsp
}
