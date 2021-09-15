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
)

type Hello struct{}

func (h *Hello) SaySrv(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) {
    rsp.Pong = "pong"
}

func (h *Hello) SayStream(c *context.Context, req *phello.SayReq) (chan *phello.SayRsp, chan error) {
    var rsp = make(chan *phello.SayRsp)
    var err = make(chan error)

    go func() {
        var count uint32
        for i := 0; i < 200; i++ {
            rsp <- &phello.SayRsp{Pong: strconv.Itoa(i)}
            atomic.AddUint32(&count, 1)
            time.Sleep(time.Second * 1)
        }

        c.Info("say stream example count is: ", atomic.LoadUint32(&count))

        close(rsp)
        close(err)
    }()

    return rsp, err
}

func (h *Hello) SayChannel(c *context.Context, req chan *phello.SayReq, errIn chan error) (
    chan *phello.SayRsp,
    chan error,
) {
    var rsp = make(chan *phello.SayRsp)
    var errs = make(chan error)

    go func() {
    QUIT:
        for {
            select {
            case _, ok := <-req:
                if !ok {
                    break QUIT
                }

                //test channel sending frequency
                time.Sleep(time.Second)
                rsp <- &phello.SayRsp{Pong: "pong"}

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
