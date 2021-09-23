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

package im

import (
    "fmt"

    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

type Im struct {
    H *Hub
    p *point
}

func (i *Im) Connect(c *context.Context, req *phello.ConnectReq, rsp *phello.ConnectRsp) {
    i.p = &point{userName: req.UserName, message: make(chan *phello.SendMessageRsp)}
    i.H.addClient(i.p)
    rsp.IsConnect = true
}

func (i *Im) Count(c *context.Context, req *phello.CountReq, rsp *phello.CountRsp) {
    rsp.Count = i.H.count()
}

func (i *Im) SendMessage(
    c *context.Context,
    req chan *phello.SendMessageReq,
    exit chan struct{},
) chan *phello.SendMessageRsp {

    var rsp = make(chan *phello.SendMessageRsp, cap(req))
    go func() {
    QUIT:
        for {
            select {
            case data, ok := <-req:
                if !ok {
                    break QUIT
                }

                i.H.broadCast <- data
                rsp <- &phello.SendMessageRsp{Message: "pong"}
            case <-exit:
                //client is closed
                //close(i.p.message)
                fmt.Println("---exit---")
            }
        }

        i.H.removeClient(i.p)
    }()

    return rsp
}
