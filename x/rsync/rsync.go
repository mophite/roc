package rsync

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"

	"github.com/go-roc/roc/internal/etcd"
)

const rsyncLockPrefix = "rocRsyncLock/"

//Acquire is a distributed lock by etcd
//try to lock with a key
//if timeout the lock will be return
//f() is the function what will be lock
//it will return a error
//key is prefix or a unique id
//tll is lock timeout setting
//tryLockTimes is backoff to retry lock
func Acquire(key string, ttl int, f func() error) error {

	if ttl <= 0 {
		ttl = 10
	}

	// get a concurrency session
	session, err := concurrency.NewSession(etcd.DefaultEtcd.Client(), concurrency.WithTTL(ttl))
	if err != nil {
		return err
	}

	defer session.Close()

	mu := concurrency.NewMutex(session, rsyncLockPrefix+key)
	err = mu.Lock(context.Background())

	//if occur a error retry lock
	if err != nil {
		return err
	}

	err = f()

	_ = mu.Unlock(context.Background())

	return err
}

func AcquireDelay(key string, ttl int, f func() error) error {

	if ttl <= 0 {
		ttl = 10
	}

	// get a concurrency session
	session, err := concurrency.NewSession(etcd.DefaultEtcd.Client(), concurrency.WithTTL(ttl))
	if err != nil {
		return err
	}

	mu := concurrency.NewMutex(session, rsyncLockPrefix+key)
	err = mu.Lock(context.Background())

	//if occur a error retry lock
	if err != nil {
		return err
	}

	//the lock is not released until the lease expires
	session.Orphan()

	err = f()

	return err
}

var ErrLock = errors.New("lock failure")

//if lock get failure during time.Second,it will be return lock err
func AcquireOnce(key string, ttl int, f func() error) error {

	if ttl <= 0 {
		ttl = 10
	}

	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// get a concurrency session
	session, err := concurrency.NewSession(
		etcd.DefaultEtcd.Client(),
		concurrency.WithContext(c),
		concurrency.WithTTL(ttl),
	)
	if err != nil {
		return ErrLock
	}

	mu := concurrency.NewMutex(session, rsyncLockPrefix+key)

	err = mu.Lock(c)

	//if occur a error retry lock
	if err != nil {
		return err
	}

	//the lock is not released until the lease expires
	session.Orphan()

	err = f()

	return err
}
