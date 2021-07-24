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

	"tutorials/proto/pim"

	"github.com/coreos/etcd/clientv3"
	im2 "tutorials/app/srv/srv.im/im"

	"github.com/go-roc/roc"
)

func main() {
	var s = roc.NewService(
		//roc.TCPAddress("127.0.0.1:8899"),
		roc.Namespace("srv.im"),
		roc.EtcdConfig(&clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		}),
	)
	pim.RegisterImServer(s, &im2.Im{H: im2.NewHub()})
	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
}
