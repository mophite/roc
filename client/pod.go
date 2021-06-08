package client

import (
	"sync"

	"roc/internal/endpoint"
	"roc/internal/transport"
)

type pod struct {
	sync.Mutex
	//service name
	name       string
	count      int
	index      uint32
	clients    []*Conn
	clientsMap map[string]*Conn
	callback   chan string
}

func newPod() *pod {
	return &pod{
		clients:    make([]*Conn, 0, 10),
		clientsMap: make(map[string]*Conn),
		callback:   make(chan string),
	}
}

func (p *pod) Add(e *endpoint.Endpoint, client transport.Client) error {
	conn, err := newConn(e, client, p.callback)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	conn.i = len(p.clients)
	p.count += 1
	p.name = e.Name
	p.clients = append(p.clients, conn)
	p.clientsMap[e.Address] = conn

	go p.watch()

	conn.working()

	return nil
}

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

func (p *pod) watch() {
	for {
		select {
		case address := <-p.callback:
			p.Del(address)
		}
	}
}
