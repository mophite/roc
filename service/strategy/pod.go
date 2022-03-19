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

package strategy

import (
	"sync"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/service/conn"
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
	clients []*conn.Conn

	//clientMap in pod
	clientsMap map[string]*conn.Conn

	clientsMapId map[string]*conn.Conn

	//when client occur a error,handler callback
	//callback is the server address
	callback chan string
}

//create a pod
func newPod() *pod {
	return &pod{
		clients:      make([]*conn.Conn, 0, 10),
		clientsMap:   make(map[string]*conn.Conn),
		clientsMapId: make(map[string]*conn.Conn),
		callback:     make(chan string),
	}
}

// Add add a client endpoint to pod
func (p *pod) Add(e *endpoint.Endpoint) error {

	p.Lock()
	defer p.Unlock()

	c, err := conn.NewConn(e, p.callback)
	if err != nil {
		return err
	}

	c.E = e

	//setting conn array cursor
	c.SetCursor(len(p.clients))

	p.count += 1
	p.serviceName = e.Name
	p.clients = append(p.clients, c)
	p.clientsMap[e.Address] = c
	p.clientsMapId[e.Id] = c

	// update callback
	go p.watch()

	// let client's conn working
	c.Working()

	return nil
}

// Del delete a client endpoint from pod
func (p *pod) Del(addr string) {
	p.Lock()
	defer p.Unlock()

	c, ok := p.clientsMap[addr]
	if ok {
		c.CloseConn()
		p.clients = append(p.clients[:c.Cursor()], p.clients[c.Cursor()+1:]...)
		delete(p.clientsMap, addr)
		delete(p.clientsMapId, c.E.Id)
		p.count -= 1
		p.index -= 1
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
