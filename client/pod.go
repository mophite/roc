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

package client

import (
	"sync"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/transport"
)

type pod struct {
	sync.Mutex

	//serviceName serviceName
	serviceName string

	//count the all clients in this pod
	count int

	//Round-robin call cursor
	index uint32

	//clients array in pod
	clients []*Conn

	//clientMap in pod
	clientsMap map[string]*Conn

	//when client occur a error,handler callback
	//callback is the server address
	callback chan string
}

//create a pod
func newPod() *pod {
	return &pod{
		clients:    make([]*Conn, 0, 10),
		clientsMap: make(map[string]*Conn),
		callback:   make(chan string),
	}
}

// Add add a client endpoint to pod
func (p *pod) Add(e *endpoint.Endpoint, client transport.Client) error {
	conn, err := newConn(e, client, p.callback)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	conn.podLength = len(p.clients)
	p.count += 1
	p.serviceName = e.Name
	p.clients = append(p.clients, conn)
	p.clientsMap[e.Address] = conn

	// watch callback
	go p.watch()

	// let client's conn working
	conn.working()

	return nil
}

// Del delete a client endpoint from pod
func (p *pod) Del(addr string) {
	p.Lock()
	defer p.Unlock()

	conn, ok := p.clientsMap[addr]
	if ok {
		conn.Close()
		p.clients = append(p.clients[:conn.Offset()], p.clients[conn.Offset():]...)
		delete(p.clientsMap, addr)
	}
}

//watch callback server address to delete
func (p *pod) watch() {
	for {
		select {
		case address := <-p.callback:
			p.Del(address)
		}
	}
}
