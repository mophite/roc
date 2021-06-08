package client

import (
	"errors"
	"sync"
	"sync/atomic"

	"roc/internal/endpoint"
	"roc/internal/namespace"
	"roc/internal/registry"
	"roc/internal/transport"
	"roc/internal/x"
	"roc/rlog"
)

var (
	ErrorNoneServer   = errors.New("server is none to use")
	ErrorNoSuchServer = errors.New("no such server")
)

type Strategy interface {
	Next(scope string) (next *Conn, err error)
	Straight(scope, address string) (next *Conn, err error)
	Close()
}

var _ Strategy = &strategy{}

type strategy struct {
	sync.Mutex
	connPerService map[string]*pod
	registry       registry.Registry
	client         transport.Client
	action         chan *registry.Action
	close          chan struct{}
}

func newStrategy(registry registry.Registry, client transport.Client) Strategy {
	s := &strategy{
		connPerService: make(map[string]*pod),
		registry:       registry,
		client:         client,
		close:          make(chan struct{}),
	}

	s.action = s.registry.Watch()

	s.lazySync()
	go s.notify()
	return s
}

func (s *strategy) getOrSet(scope string) (*pod, error) {
	p, ok := s.connPerService[scope]
	if !ok || p.count == 0 {
		s.Lock()
		defer s.Unlock()

		e, err := s.registry.Next(scope)
		if err != nil {
			return nil, err
		}
		err = s.sync(e)
		if err != nil {
			return nil, err
		}
		return s.connPerService[e.Scope], nil
	}
	//pod must available
	return p, nil
}

func (s *strategy) Next(scope string) (next *Conn, err error) {
	p, err := s.getOrSet(scope)
	if err != nil {
		return nil, err
	}

	var conn *Conn
	for i := 0; i < p.count; i++ {
		conn = p.clients[(int(atomic.AddUint32(&p.index, 1))-1)%p.count]
		if conn.isWorking() {
			break
		}
	}

	if conn == nil || !conn.isWorking() {
		return nil, ErrorNoneServer
	}

	return conn, nil
}

func (s *strategy) Straight(scope, address string) (next *Conn, err error) {
	p, err := s.getOrSet(scope)
	if err != nil {
		return nil, err
	}

	conn, ok := p.clientsMap[address]
	if !ok || !conn.isWorking() {
		return nil, ErrorNoSuchServer
	}

	return conn, nil
}

func (s *strategy) lazySync() {
	s.Lock()
	defer s.Unlock()

	es, err := s.registry.List()
	if err != nil {
		rlog.Debug(err)
		return
	}

	for _, e := range es {
		s.sync(e)
	}
}

func (s *strategy) sync(e *endpoint.Endpoint) error {
	p, ok := s.connPerService[e.Scope]
	if !ok {
		p = newPod()
	}

	err := p.Add(e, s.client)
	if err != nil {
		rlog.Error(err)
		return err
	}

	s.connPerService[e.Scope] = p

	return nil
}

func (s *strategy) notify() {

QUIT:
	for {
		select {
		case act := <-s.action:
			s.watch(act)
			rlog.Debug("watch endpoint.Endpoint was changed", x.MustMarshalString(act))
		case <-s.close:
			break QUIT
		}
	}
}

func (s *strategy) watch(act *registry.Action) {
	s.Lock()
	defer s.Unlock()

	switch act.Act {
	case namespace.WatcherCreate:
		_ = s.sync(act.E)
	case namespace.WatcherDelete:
		p, ok := s.connPerService[act.E.Scope]
		if ok {
			p.Del(act.E.Address)
		}
	}
}

func (s *strategy) Close() {
	for _, p := range s.connPerService {
		for _, client := range p.clientsMap {
			client.Close()
		}
	}

	s.close <- struct{}{}
}
