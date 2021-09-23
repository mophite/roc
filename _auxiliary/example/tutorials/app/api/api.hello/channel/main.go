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

package imChannel

import (
    "fmt"
    "strconv"
    "sync/atomic"

    "github.com/go-roc/roc/_auxiliary/example/tutorials/internal/ipc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

func Channel() {

    var req = make(chan *phello.SayReq, 100)
    go func() {
        for i := 0; i < 3; i++ {

            //test sending frequency
            //time.Sleep(time.Second)
            req <- &phello.SayReq{Ping: strconv.Itoa(i)}

            //if i == 20 {
            //	errsIn <- errors.New("send a test error")
            //	break
            //}
        }

        close(req)
    }()

    rsp, exit := ipc.SayChannel(context.Background(), req)

    var count uint32

    var done = make(chan struct{})
    go func() {
    QUIT:
        for {
            select {
            case b, ok := <-rsp:
                if ok {
                    fmt.Println("------receive from srv.hello----", b.Pong)
                    atomic.AddUint32(&count, 1)
                } else {
                    break QUIT
                }

                if atomic.LoadUint32(&count) == 3 {
                    break QUIT
                }

            case <-exit:
                break QUIT
            }
        }
        done <- struct{}{}

        fmt.Println("say handler count is: ", atomic.LoadUint32(&count))
    }()

    <-done
}
