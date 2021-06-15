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

package main

import (
	"fmt"
	"roc/_auxiliary/example/tutorials/proto/pbhello"
	"roc/client"
	"roc/parcel/context"
	"sync/atomic"
)

var helloClient = pbhello.NewHelloWorldClient(client.NewRocClient())
var opt = client.WithName("srv.hello")

func main() {

	var req = make(chan *pbhello.SayReq, 100)
	var errsIn = make(chan error)
	go func() {
		for i := 0; i < 50; i++ {

			//test sending frequency
			//time.Sleep(time.Second)
			req <- &pbhello.SayReq{Inc: uint32(i)}

			//if i == 20 {
			//	errsIn <- errors.New("send a test error")
			//	break
			//}
		}

		close(req)
		close(errsIn)
	}()

	rsp, errs := helloClient.SayChannel(context.Background(), req, errsIn, opt)

	var count uint32

	var done = make(chan struct{})
	go func() {
		var err error
	QUIT:
		for {
			select {
			case b, ok := <-rsp:
				if ok {
					fmt.Println("------receive from srv.hello----", b.Inc)
					atomic.AddUint32(&count, 1)
				} else {
					break QUIT
				}
			case err = <-errs:
				if err != nil {
					break QUIT
				}
			}
		}
		done <- struct{}{}

		fmt.Println("say handler count is: ", atomic.LoadUint32(&count))
	}()

	<-done
}
