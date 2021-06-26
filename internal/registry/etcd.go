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

package registry

import (
	"github.com/coreos/etcd/clientv3"

	"github.com/go-roc/roc/internal/etcd"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/x"
	"github.com/go-roc/roc/rlog"
)

//etcd implementation of service discovery
type etcdRegistry struct {

	//registry option
	opts Option

	//etcd instance
	e *etcd.Etcd

	//watch instance
	watch *etcd.Watch
}

// NewRegistry create a new registry with etcd
func NewRegistry(opts ...Options) Registry {
	r := &etcdRegistry{opts: newOpts(opts...)}

	r.e = etcd.DefaultEtcd
	return r
}

// Register register one endpoint to etcd
func (s *etcdRegistry) Register(e *endpoint.Endpoint) error {
	return s.e.PutWithLease(e.Absolute, x.MustMarshalString(e))
}

// Next return a endpoint
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

// List get all endpoint from etcd
func (s *etcdRegistry) List() (services []*endpoint.Endpoint, err error) {
	b, err := s.e.GetWithList(namespace.SplicingPrefix(s.opts.schema, ""), clientv3.WithPrefix())
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

// Deregister deregister a endpoint ,remove it from etcd
func (s *etcdRegistry) Deregister(e *endpoint.Endpoint) error {
	return s.e.Delete(e.Absolute)
}

func (s *etcdRegistry) Name() string {
	return "etcd"
}

func (s *etcdRegistry) Watch() chan *Action {
	var r = make(chan *Action)

	s.watch = etcd.NewEtcdWatch(s.opts.schema, s.e.Client())

	go func() {
		for v := range s.watch.Watch(s.opts.schema) {
			for key, value := range v.B {
				var e endpoint.Endpoint
				err := x.Jsoniter.Unmarshal(value, &e)
				if err != nil {
					rlog.Warnf("action=%s |err=%v", x.MustMarshalString(v.B), err)
					continue
				}
				r <- &Action{
					Act: v.Act,
					E:   &e,
					Key: key,
				}
			}
		}

		close(r)
	}()

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
