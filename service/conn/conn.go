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

package conn

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/transport"
	rs "github.com/go-roc/roc/internal/transport/rsocket"
)

// state is mark conn state,conn must safe
type state = uint32

const (
	// StateBlock block is unavailable state
	StateBlock state = 0x01 + iota
	// StateReady ready is unavailable state
	StateReady

	// StateWorking is available state
	StateWorking

	// errCountDelta is record the number of connection failures
	errCountDelta = 3
)

// Conn include transport client
//
type Conn struct {
	sync.Mutex

	cursor int

	// conn state
	state state

	// current conn occur error count
	errCount uint32

	// client object
	client transport.Client
}

func (c *Conn) SetCursor(i int) {
	c.cursor = i
}

func (c *Conn) Cursor() int {
	return c.cursor
}

// GrowErrorCount error safe grow one
func (c *Conn) GrowErrorCount() uint32 {
	return atomic.AddUint32(&c.errCount, 1)
}

// Working swap state to working
func (c *Conn) Working() {
	atomic.SwapUint32(&c.state, StateWorking)
}

// Block swap state to block
func (c *Conn) Block() {
	atomic.SwapUint32(&c.state, StateBlock)
}

// Ready swap state to ready
func (c *Conn) Ready() {
	atomic.SwapUint32(&c.state, StateReady)
}

// GetState get state
func (c *Conn) GetState() state {
	return atomic.LoadUint32(&c.state)
}

// IsWorking judge state is working
func (c *Conn) IsWorking() bool {
	return c.GetState() == StateWorking
}

// IsBlock judge state is block
func (c *Conn) IsBlock() bool {
	return c.GetState() == StateBlock
}

// GrowError grow error and let the error conn retry working util conn is out of serviceName
func (c *Conn) GrowError() {
	c.Lock()
	defer c.Unlock()

	if c.GrowErrorCount() > errCountDelta && c.IsWorking() {
		// let conn block
		c.Block()
		go func() {
			select {
			case <-time.After(time.Second * 3):
				// let conn working
				// if conn is out of serviceName,this is not effect
				//todo try to ping ,if ok let client working
				//if close ,don't do anything
				c.Working()
			}
		}()
	}
}

// NewConn is create a conn
// closeCallBack is the conn client occur error and callback
func NewConn(
	connectTimeout,
	keepaliveInterval,
	keepaliveLifetime time.Duration,
	e *endpoint.Endpoint,
	closeCallback chan string,
) (*Conn, error) {
	client := rs.NewClient(
		connectTimeout,
		keepaliveInterval,
		keepaliveLifetime,
	)
	err := client.Dial(e, closeCallback)
	if err != nil {
		return nil, err
	}

	c := &Conn{client: client}

	// change state to ready
	c.Ready()

	return c, nil
}

// Client get client
func (c *Conn) Client() transport.Client {
	return c.client
}

// Close close client
func (c *Conn) Close() {
	c.Block()

	c.Client().Close()

	c.client = nil
}
