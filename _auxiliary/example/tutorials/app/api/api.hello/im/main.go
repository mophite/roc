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

package ImTest

import (
    "fmt"
    "strconv"
    "sync/atomic"
    "time"

    "github.com/go-roc/roc/_auxiliary/example/tutorials/internal/ipc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/rlog"
)

func Im() {

    cRsp, err := ipc.Connect(context.Background(), &phello.ConnectReq{UserName: "roc"})
    if err != nil {
        rlog.Error(err)
        return
    }

    if !cRsp.IsConnect {
        return
    }

    var req = make(chan *phello.SendMessageReq)
    go func() {
        for i := 0; i < 3; i++ {
            req <- &phello.SendMessageReq{Message: "im - " + strconv.Itoa(i)}
        }

        //must be closed
        time.Sleep(time.Second)
        close(req)
    }()
    rsp := ipc.SendMessage(context.Background(), req)

    var count uint32

    var done = make(chan struct{})
    go func() {
    QUIT:
        for {
            select {
            case b, ok := <-rsp:
                if ok {
                    fmt.Println("------receive from srv.im----", b.Message)
                    atomic.AddUint32(&count, 1)
                } else {
                    break QUIT
                }
            }
        }

        done <- struct{}{}

        fmt.Println("say handler count is: ", atomic.LoadUint32(&count))
    }()
    <-done
}
