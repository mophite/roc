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

package etcd

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"

	"roc/internal/x"
	"roc/rlog"
)

type Etcd struct {
	lock sync.RWMutex

	//etcd client v3
	client *clientv3.Client

	//etcd leaseId
	leaseId clientv3.LeaseID

	//use leaseId to keepalive
	leaseKeepaliveChan chan *clientv3.LeaseKeepAliveResponse

	//etcd config
	config *clientv3.Config

	//timeout setting
	timeout time.Duration

	//leaseTLL setting
	leaseTLL int64
}

// NewEtcd create etcd
// if config is nil,use default config setting
func NewEtcd(timeout time.Duration, leaseTLL int64, config *clientv3.Config) *Etcd {
	s := &Etcd{
		leaseKeepaliveChan: make(chan *clientv3.LeaseKeepAliveResponse),
		config:             config,
		timeout:            timeout,
		leaseTLL:           leaseTLL,
	}

	if s.config == nil {
		s.config = &clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		}
	}

	var err error
	s.client, err = clientv3.New(*s.config)
	if err != nil {
		panic(err)
	}
	return s
}

// Client get etcd client
func (s *Etcd) Client() *clientv3.Client {
	return s.client
}

// PutWithLease put one key/value to etcd with lease setting
func (s *Etcd) PutWithLease(key, value string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.timeout)
	defer cancel()

	rsp, err := clientv3.NewLease(s.client).Grant(ctx, s.leaseTLL)
	if err != nil {
		return err
	}

	s.leaseId = rsp.ID

	ch, err := s.client.KeepAlive(context.TODO(), rsp.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case c := <-s.leaseKeepaliveChan: // if leaseKeepaliveChan is nil,lease keepalive stop!
				if c == nil {
					rlog.Warnf("etcd leaseKeepalive stop! leaseID: %d prefix:%s value:%s", s.leaseId, key, value)
					return
				}
			}
		}
	}()

	s.leaseKeepaliveChan <- <-ch
	_, err = s.client.Put(context.TODO(), key, value, clientv3.WithLease(s.leaseId))
	if err != nil {
		switch err {
		case context.Canceled:
			rlog.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			rlog.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			rlog.Fatalf("client-side error: %v", err)
		default:
			rlog.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	return err
}

// Put put one key/value to etcd with no lease setting
func (s *Etcd) Put(key, value string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := s.client.Put(context.Background(), key, value)
	if err != nil {
		switch err {
		case context.Canceled:
			rlog.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			rlog.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			rlog.Fatalf("client-side error: %v", err)
		default:
			rlog.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}

	return err
}

// GetWithLastKey get value with last key
func (s *Etcd) GetWithLastKey(key string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	rsp, err := s.client.Get(ctx, key, clientv3.WithLastKey()...)
	if err != nil {
		return nil, err
	}

	if rsp.Count < 1 {
		return nil, errors.New("GetWithLastKey is none by etcd")
	}

	return rsp.Kvs[int(rsp.Count-1)].Value, nil
}

// GetWithKey get value with key
func (s *Etcd) GetWithKey(key string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	rsp, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if rsp.Count < 1 {
		return nil, errors.New("GetWithKey is none by etcd")
	}

	return rsp.Kvs[int(rsp.Count-1)].Value, nil
}

func (s *Etcd) GetWithList(key string) (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	rsp, err := s.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	if rsp.Count < 1 {
		return nil, errors.New("GetWithList is none by etcd")
	}

	var r = make(map[string][]byte, rsp.Count)

	for i := range rsp.Kvs {
		r[x.BytesToString(rsp.Kvs[i].Key)] = rsp.Kvs[i].Value
	}

	return r, nil
}

func (s *Etcd) Delete(key string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ctx, cancel := context.WithTimeout(context.TODO(), s.timeout)
	defer cancel()

	_, err := s.client.Delete(ctx, key)

	return err
}

//revoke lease
func (s *Etcd) revoke() error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.timeout)
	defer cancel()
	_, err := s.client.Revoke(ctx, s.leaseId)
	return err
}

func (s *Etcd) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.client != nil {
		if s.leaseId > 0 {
			err := s.revoke()
			if err != nil {
				rlog.Error(err)
			}
		}
		err := s.client.Close()
		if err != nil {
			rlog.Error(err)
		}
	}

	if s.leaseKeepaliveChan != nil {
		close(s.leaseKeepaliveChan)
	}
}
