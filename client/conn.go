package client

import (
	"sync"
	"sync/atomic"
	"time"

	"roc/internal/endpoint"
	"roc/internal/transport"
)

type state = uint32

const (
	StateBlock state = 0x01 + iota
	StateReady
	StateWorking

	errCountDelta = 3
)

type Conn struct {
	sync.Mutex
	i        int
	state    state
	errCount uint32
	client   transport.Client
}

func (c *Conn) growErrorCount() uint32 {
	return atomic.AddUint32(&c.errCount, 1)
}

func (c *Conn) working() {
	atomic.SwapUint32(&c.state, StateWorking)
}

func (c *Conn) block() {
	atomic.SwapUint32(&c.state, StateBlock)
}

func (c *Conn) ready() {
	atomic.SwapUint32(&c.state, StateReady)
}

func (c *Conn) getState() state {
	return atomic.LoadUint32(&c.state)
}

func (c *Conn) isWorking() bool {
	return c.getState() == StateWorking
}

func (c *Conn) isBlock() bool {
	return c.getState() == StateBlock
}

func (c *Conn) growError() {
	c.Lock()
	defer c.Unlock()

	if c.growErrorCount() > errCountDelta && !c.isBlock() {
		c.block()
		go func() {
			select {
			case <-time.After(time.Second * 3):
				c.working()
			}
		}()
	}
}

func newConn(e *endpoint.Endpoint, client transport.Client, ch chan string) (*Conn, error) {
	err := client.Dial(e, ch)
	if err != nil {
		return nil, err
	}

	c := &Conn{client: client}

	c.ready()

	return c, nil
}

func (c *Conn) Offset() int {
	return c.i
}

func (c *Conn) Client() transport.Client {
	return c.client
}

func (c *Conn) Close() {
	c.block()
	c.Client().Close()
	c.client = nil
}
