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
	"roc/_auxiliary/example/tutorials/proto/pbim"
	"roc/client"
	"roc/parcel/context"
	"strconv"
	"sync/atomic"
)

var imClient = pbim.NewImClient(client.NewRocClient())
var opt = client.WithName("srv.im")

func main() {

	cRsp, err := imClient.Connect(context.Background(), &pbim.ConnectReq{UserName: "roc"})
	if err != nil {
		panic(err)
	}

	if !cRsp.IsConnect {
		return
	}

	var req = make(chan *pbim.SendMessageReq)
	var errsIn = make(chan error)
	go func() {
		for i := 0; i < 3; i++ {
			req <- &pbim.SendMessageReq{Message: "im - " + strconv.Itoa(i)}
		}

		//close(req)
		//close(errsIn)
	}()

	rsp, errs := imClient.SendMessage(context.Background(), req, errsIn, opt)

	var count uint32

	var done = make(chan struct{})
	go func() {
		var err error
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
