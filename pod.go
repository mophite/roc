package roc

import (
	"fmt"
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

	p.Lock()
	defer p.Unlock()

	conn, err := newConn(e, client, p.callback)
	if err != nil {
		return err
	}

	//setting conn array cursor
	conn.cursor = len(p.clients)

	fmt.Println("----1--", conn.cursor, e.Name, e.Address)
	p.count += 1
	p.serviceName = e.Name
	p.clients = append(p.clients, conn)
	p.clientsMap[e.Address] = conn
	fmt.Println("-----2--", p.clients)

	// update callback
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
		p.clients = append(p.clients[:conn.cursor], p.clients[conn.cursor+1:]...)
		delete(p.clientsMap, addr)
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
