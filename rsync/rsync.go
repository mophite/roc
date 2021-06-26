package rsync

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"

	"github.com/go-roc/roc/internal/backoff"
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
func Acquire(key string, ttl, tryLockTimes int, f func() error) error {
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
	err = mu.Lock(context.TODO())

	//if occur a error retry lock
	if err != nil {
		bf := backoff.NewBackoff()
		for i := 0; i < tryLockTimes; i++ {
			time.Sleep(bf.Next(i))
			err = mu.Lock(context.TODO())
			if err != nil {
				continue
			}
			break
		}

		if err != nil {
			return err
		}
	}

	err = f()

	_ = mu.Unlock(context.TODO())

	return err
}
