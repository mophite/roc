package roc

import (
    "errors"
    "sync"
    "sync/atomic"

    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/internal/registry"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/x"
)

var (
    ErrorNoneServer   = errors.New("server is none to use")
    ErrorNoSuchServer = errors.New("no such server")
)

type Strategy interface {

    //Next Round-robin scheduling
    Next(scope string) (next *Conn, err error)

    //Straight direct call
    Straight(scope, address string) (next *Conn, err error)

    //Close Strategy
    Close()
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

    //filter local service registry
    localEndpoint *endpoint.Endpoint

    service *Service
}

// create a strategy
func newStrategy(
    local *endpoint.Endpoint,
    registry registry.Registry,
    service *Service,
) Strategy {
    s := &strategy{
        connPerService: make(map[string]*pod),
        registry:       registry,
        close:          make(chan struct{}),
        localEndpoint:  local,
        service:        service,
    }

    //receive registry update notify
    s.action = s.registry.Watch()

    //Synchronize all existing services
    s.lazySync()

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

        return s.connPerService[scope], nil
    }
    //pod must available
    return p, nil
}

// Next Round-robin next
func (s *strategy) Next(scope string) (*Conn, error) {
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

// Straight direct invoke
func (s *strategy) Straight(scope, address string) (*Conn, error) {
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
    if e.Scope == s.localEndpoint.Scope {
        return nil
    }

    p, ok := s.connPerService[e.Scope]
    if !ok {
        p = newPod()
    }

    err := p.Add(e, s.service)
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

// Close strategy close
func (s *strategy) Close() {
    for _, p := range s.connPerService {
        for _, client := range p.clientsMap {
            client.Close()
        }
    }

    s.close <- struct{}{}
}
