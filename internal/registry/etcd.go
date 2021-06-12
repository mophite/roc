package registry

import (
	"go.etcd.io/etcd/clientv3"
	"roc/internal/endpoint"
	"roc/internal/etcd"
	"roc/internal/namespace"
	"roc/internal/x"
	"roc/rlog"
)

type etcdRegistry struct {
	opts  Option
	e     *etcd.Etcd
	watch *etcd.Watch
}

func NewRegistry(opts ...Options) Registry {
	r := &etcdRegistry{opts: newOpts(opts...)}

	if r.opts.etcdConfig == nil {
		r.opts.etcdConfig = new(clientv3.Config)
	}

	if len(r.opts.address) > 0 {
		r.opts.etcdConfig.Endpoints = r.opts.address
	}

	r.e = etcd.NewEtcd(r.opts.timeout, r.opts.leaseTLL, r.opts.etcdConfig)
	return r
}

func (s *etcdRegistry) Register(e *endpoint.Endpoint) error {
	return s.e.PutWithLease(e.Absolute, x.MustMarshalString(e))
}

func (s *etcdRegistry) Next(scope string) (*endpoint.Endpoint, error) {

	b, err := s.e.GetWithLastKey(namespace.SplicingPrefix(s.opts.schema, scope))
	if err != nil {
		return nil, err
	}

	var e endpoint.Endpoint

	err = x.Jsoniter.Unmarshal(b, &e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *etcdRegistry) List() (services []*endpoint.Endpoint, err error) {
	b, err := s.e.GetWithList(namespace.SplicingPrefix(s.opts.schema, ""))
	if err != nil {
		return nil, err
	}

	for _, v := range b {
		var e endpoint.Endpoint
		err = x.Jsoniter.Unmarshal(v, &e)
		if err != nil {
			continue
		}
		services = append(services, &e)
	}
	return
}

func (s *etcdRegistry) Deregister(e *endpoint.Endpoint) error {
	return s.e.Delete(e.Absolute)
}

func (s *etcdRegistry) Name() string {
	return "etcd"
}

func (s *etcdRegistry) Watch() chan *Action {
	var r = make(chan *Action)
	go func() {
		for v := range s.watch.Watch(s.opts.schema) {
			for _, value := range v.B {
				var e endpoint.Endpoint
				err := x.Jsoniter.Unmarshal(value, &e)
				if err != nil {
					rlog.Warnf("action=%s |err=%v", x.MustMarshalString(v.B), err)
					continue
				}
				r <- &Action{
					Act: v.Act,
					E:   &e,
				}
			}
		}

		close(r)
	}()

	s.watch = etcd.NewEtcdWatch(s.opts.schema, s.e.Client())
	return r
}

func (s *etcdRegistry) Close() {
	if s.e != nil {
		s.e.Close()
	}

	if s.watch != nil {
		s.watch.Close()
	}
}
