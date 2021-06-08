package registry

import (
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

func DefaultEtcdRegistry(opts ...Options) Registry {
	r := &etcdRegistry{opts: newOpts(opts...)}
	r.e = etcd.NewEtcd(r.opts.timeout, r.opts.leaseTLL, r.opts.etcdConfig)
	r.watch = etcd.NewEtcdWatch(r.opts.schema, r.e.Client())
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

	return r
}

func (s *etcdRegistry) Close() {
	s.e.Close()
	s.watch.Close()
}
