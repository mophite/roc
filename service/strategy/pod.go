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

	//when client occur a error,handler callback
	//callback is the server address
	callback chan string
}

//create a pod
func newPod() *pod {
	return &pod{
		clients:    make([]*conn.Conn, 0, 10),
		clientsMap: make(map[string]*conn.Conn),
		callback:   make(chan string),
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

	//setting conn array cursor
	c.SetCursor(len(p.clients))

	p.count += 1
	p.serviceName = e.Name
	p.clients = append(p.clients, c)
	p.clientsMap[e.Address] = c

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
		c.Close()
		p.clients = append(p.clients[:c.Cursor()], p.clients[c.Cursor()+1:]...)
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
