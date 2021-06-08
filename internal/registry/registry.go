package registry

import (
	"roc/internal/endpoint"
	"roc/internal/namespace"
)

type Registry interface {
	Watcher
	Register(e *endpoint.Endpoint) error
	Deregister(e *endpoint.Endpoint) error
	Next(scope string) (*endpoint.Endpoint, error)
	List() ([]*endpoint.Endpoint, error)
	Name() string
	Close()
}

type Watcher interface {
	Watch() chan *Action
}

type Action struct {
	Act namespace.WatcherAction
	E   *endpoint.Endpoint
}
