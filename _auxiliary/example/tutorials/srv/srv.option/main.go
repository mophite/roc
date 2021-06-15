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
	"roc/_auxiliary/example/tutorials/srv/srv.hello/hello"
	"roc/internal/registry"
	"roc/server"
)

func main() {
	var opt = registry.Address([]string{"82.157.14.79:2379"})

	var s = server.NewRocServer(
		server.Namespace("srv.hello"),
		server.TCPAddress("127.0.0.1:8089"),
		server.Registry(registry.NewRegistry(opt)),
	)
	pbhello.RegisterHelloWorldServer(s, &hello.Hello{})
	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
}
