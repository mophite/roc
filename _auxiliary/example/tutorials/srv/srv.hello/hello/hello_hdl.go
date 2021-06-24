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
	"sync/atomic"
	"time"

	"github.com/go-roc/roc"
	"github.com/go-roc/roc/_auxiliary/example/tutorials/proto/pbhello"
	"github.com/go-roc/roc/parcel/context"
)

type Hello struct {
	Service *roc.Service
}

func (h *Hello) SayStream(c *context.Context, req *pbhello.SayReq) (chan *pbhello.SayRsp, chan error) {
	var rsp = make(chan *pbhello.SayRsp)
	var err = make(chan error)

	go func() {
		var count uint32
		for i := 0; i < 200; i++ {
			rsp <- &pbhello.SayRsp{Inc: req.Inc + uint32(i)}
			atomic.AddUint32(&count, 1)
			time.Sleep(time.Second * 1)
		}

		c.Info("say stream example count is: ", atomic.LoadUint32(&count))

		close(rsp)
		close(err)
	}()

	return rsp, err
}

func (h *Hello) SayChannel(c *context.Context, req chan *pbhello.SayReq, errIn chan error) (chan *pbhello.SayRsp, chan error) {
	var rsp = make(chan *pbhello.SayRsp)
	var errs = make(chan error)

	go func() {
	QUIT:
		for {
			select {
			case data, ok := <-req:
				if !ok {
					break QUIT
				}

				//test channel sending frequency
				time.Sleep(time.Second)
				rsp <- &pbhello.SayRsp{Inc: data.Inc + uint32(1)}

			case e := <-errIn:
				if e != nil {
					errs <- e
				}
			}
		}

		close(rsp)
		close(errs)
	}()

	return rsp, errs
}

func (h *Hello) Say(c *context.Context, req *pbhello.SayReq) (rsp *pbhello.SayRsp, err error) {
	//  when set timeout is time.Second*1,it's will occur cancelled error
	time.Sleep(time.Second * 2)
	return &pbhello.SayRsp{Inc: req.Inc + 1}, nil
}
