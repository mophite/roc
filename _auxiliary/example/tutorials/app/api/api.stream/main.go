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
	"sync/atomic"

	"tutorials/proto/phello"

	"github.com/coreos/etcd/clientv3"

	"github.com/go-roc/roc"
)

var helloClient = phello.NewHelloWorldClient(
	roc.NewService(
		roc.TCPAddress("127.0.0.1:8899"),
		roc.Namespace("srv.im"),
		roc.EtcdConfig(
			&clientv3.Config{
				Endpoints: []string{"127.0.0.1:2379"},
			},
		),
	),
)
var opt = roc.WithName("srv.hello")

func main() {

	rsp, errs := helloClient.SayStream(context.Background(), &phello.SayReq{Inc: 1}, opt)

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
