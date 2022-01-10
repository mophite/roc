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
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/rlog"
	"github.com/go-roc/roc/service/conn"
	"github.com/go-roc/roc/x"
)

var (
	ErrorNoneServer   = errors.New("server is none to use")
	ErrorNoSuchServer = errors.New("no such server")
)

type Strategy interface {

	//Next Round-robin scheduling
	Next(scope string) (next *conn.Conn, err error)

	//Straight direct call
	Straight(scope, address string) (next *conn.Conn, err error)

	//CloseStrategy Strategy
	CloseStrategy()
}

var _ Strategy = &strategy{}

type strategy struct {
	sync.Mutex

	//per service & multiple conn
	connPerService map[string]*pod

	//discover registry
	registry registry.Registry

	//registry update callback action
	action chan *registry.Action

	//close strategy signal
	close chan struct{}

	localEndpoint *endpoint.Endpoint
}

// NewStrategy create a strategy
func NewStrategy(
	local *endpoint.Endpoint,
	registry registry.Registry,
) Strategy {
	s := &strategy{
		connPerService: make(map[string]*pod),
		registry:       registry,
		close:          make(chan struct{}),
		localEndpoint:  local,
	}

	//receive registry update notify
	s.action = s.registry.Watch()

	//Synchronize all existing services
	//s.lazySync()

	//handler registry notify
	go s.notify()

	return s
}

//get a pod,if is nil ,create a new pod
func (s *strategy) getOrSet(scope string) (*pod, error) {
	p, ok := s.connPerService[scope]
	if !ok || p.count == 0 {
		s.Lock()
		defer s.Unlock()

		e, err := s.registry.List()
		if err != nil {
			return nil, err
		}

		for i := range e {
			if e[i].Scope == scope {
				err = s.sync(e[i])
				if err != nil {
					return nil, err
				}
			}
		}

		v, ok := s.connPerService[scope]
		if !ok {
			return nil, fmt.Errorf("no such scope node service [%s]", scope)
		}
		return v, nil
	}
	//pod must available
	return p, nil
}

// Next Round-robin next
func (s *strategy) Next(scope string) (*conn.Conn, error) {
	p, err := s.getOrSet(scope)
	if err != nil {
		return nil, err
	}

	var c *conn.Conn
	for i := 0; i < p.count; i++ {
		c = p.clients[(int(atomic.AddUint32(&p.index, 1))-1)%p.count]
		if c.IsWorking() {
			break
		}
	}

	if c == nil || !c.IsWorking() {
		return nil, ErrorNoneServer
	}

	return c, nil
}

// Straight direct invoke
func (s *strategy) Straight(scope, address string) (*conn.Conn, error) {
	p, err := s.getOrSet(scope)
	if err != nil {
		return nil, err
	}

	c, ok := p.clientsMap[address]
	if !ok || !c.IsWorking() {
		return nil, ErrorNoSuchServer
	}

	return c, nil
}

//Synchronize all existing services
func (s *strategy) lazySync() {
	s.Lock()
	defer s.Unlock()

	es, err := s.registry.List()
	if err != nil {
		rlog.Error(err)
		return
	}

	for _, e := range es {
		_ = s.sync(e)
	}
}

//Synchronize one services
func (s *strategy) sync(e *endpoint.Endpoint) error {

	//filter local service registry
	if e.Name == s.localEndpoint.Name {
		return nil
	}

	//if reflect.DeepEqual(e, s.localEndpoint) {
	//	return nil
	//}

	p, ok := s.connPerService[e.Scope]
	if !ok {
		p = newPod()
	}

	err := p.Add(e)
	if err != nil {
		rlog.Error(err)
		return err
	}

	s.connPerService[e.Scope] = p

	return nil
}

//receive a registry notify callback
func (s *strategy) notify() {

QUIT:
	for {
		select {
		case act := <-s.action:
			rlog.Debug("update endpoint.Endpoint was changed", x.MustMarshalString(act))

			s.update(act)
		case <-s.close:
			break QUIT
		}
	}
}

//update a registry notify
func (s *strategy) update(act *registry.Action) {
	s.Lock()
	defer s.Unlock()

	switch act.Act {
	case namespace.WatcherCreate, namespace.WatcherUpdate:
		_ = s.sync(act.E)
	case namespace.WatcherDelete:
		p, ok := s.connPerService[act.E.Scope]
		if ok {
			p.Del(act.E.Address)
		}
	}
}

// CloseStrategy strategy close
func (s *strategy) CloseStrategy() {
	for _, p := range s.connPerService {
		for _, client := range p.clientsMap {
			client.CloseConn()
		}
	}

	s.close <- struct{}{}
}
