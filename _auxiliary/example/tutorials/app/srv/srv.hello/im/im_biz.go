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

package im

import (
    "sync"

    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/gogo/protobuf/proto"
)

func NewHub() *Hub {
    h := &Hub{
        lock:         new(sync.RWMutex),
        connectCount: 0,
        clients:      make(map[string]*point),
        broadCast:    make(chan *phello.SendMessageReq),
    }
    go h.poller()
    return h
}

type Hub struct {
    lock         *sync.RWMutex
    connectCount uint32
    clients      map[string]*point
    broadCast    chan *phello.SendMessageReq
}

type point struct {
    userName string
    message  chan proto.Message
}

func (h *Hub) count() uint32 {
    return h.connectCount
}

func (h *Hub) addClient(p *point) {
    if _, ok := h.clients[p.userName]; !ok {
        h.lock.RLock()
        h.clients[p.userName] = p
        h.connectCount += 1
        h.lock.RUnlock()
    }
}

func (h *Hub) removeClient(p *point) {
    if _, ok := h.clients[p.userName]; ok {
        h.lock.RLock()
        delete(h.clients, p.userName)
        h.connectCount -= 1
        h.lock.RUnlock()
    }
}

func (h *Hub) poller() {

    for {
        select {
        case b := <-h.broadCast:
            go func() {
                for userName, _ := range h.clients {
                    h.clients[userName].message <- &phello.SendMessageRsp{Message: b.Message}
                }
            }()

            // todo some thing
        }
    }
}
