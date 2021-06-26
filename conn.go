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

package roc

import (
    "sync"
    "sync/atomic"
    "time"

    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/transport"
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

// error safe grow one
func (c *Conn) growErrorCount() uint32 {
    return atomic.AddUint32(&c.errCount, 1)
}

// swap state to working
func (c *Conn) working() {
    atomic.SwapUint32(&c.state, StateWorking)
}

// swap state to block
func (c *Conn) block() {
    atomic.SwapUint32(&c.state, StateBlock)
}

// swap state to ready
func (c *Conn) ready() {
    atomic.SwapUint32(&c.state, StateReady)
}

// get state
func (c *Conn) getState() state {
    return atomic.LoadUint32(&c.state)
}

// judge state is working
func (c *Conn) isWorking() bool {
    return c.getState() == StateWorking
}

// judge state is block
func (c *Conn) isBlock() bool {
    return c.getState() == StateBlock
}

// grow error and let the error conn retry working util conn is out of serviceName
func (c *Conn) growError() {
    c.Lock()
    defer c.Unlock()

    if c.growErrorCount() > errCountDelta && !c.isBlock() {
        // let conn block
        c.block()
        go func() {
            select {
            case <-time.After(time.Second * 3):
                // let conn working
                // if conn is out of serviceName,this is not effect
                c.working()
            }
        }()
    }
}

// newConn is create a conn
// closeCallBack is the conn client occur error and callback
func newConn(e *endpoint.Endpoint, client transport.Client, closeCallback chan string) (*Conn, error) {
    err := client.Dial(e, closeCallback)
    if err != nil {
        return nil, err
    }

    c := &Conn{client: client}

    // change state to ready
    c.ready()

    return c, nil
}

// Client get client
func (c *Conn) Client() transport.Client {
    return c.client
}

// Close close client
func (c *Conn) Close() {
    c.block()

    //client.Close will be close wrong connection what you don't want
    //because rsocket is duplex conn
    //c.Client().Close()
    c.block()
    c.client = nil
}
