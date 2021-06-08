package etcd

import (
	"bytes"
	"context"
	"go.etcd.io/etcd/clientv3"
	"roc/internal/x"
	"strings"

	"roc/internal/namespace"
	"roc/rlog"
)

type Action struct {
	Act namespace.WatcherAction
	B   map[string][]byte
}

type Watch struct {
	exit   chan struct{}
	client *clientv3.Client
	wc     clientv3.WatchChan
}

func NewEtcdWatch(prefix string, client *clientv3.Client) *Watch {

	var w = &Watch{
		exit:   make(chan struct{}),
		client: client,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-w.exit
		cancel()
	}()

	w.wc = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())

	return w
}

func (w *Watch) Watch(prefix string) chan *Action {
	c := make(chan *Action)

	go func() {
		for v := range w.wc {
			if v.Err() != nil {
				rlog.Error("etcd watch err ", v.Err())
				continue
			}

			for _, event := range v.Events {

				if !strings.Contains(x.BytesToString(event.Kv.Key), prefix) {
					continue
				}

				var (
					a   namespace.WatcherAction
					key string
					b   = new(bytes.Buffer)
				)

				switch event.Type {
				case clientv3.EventTypePut:

					a = namespace.WatcherCreate

					if event.IsModify() {
						a = namespace.WatcherUpdate
					}

					if event.IsCreate() {
						a = namespace.WatcherCreate
					}

					key = x.BytesToString(event.Kv.Key)
					b.Write(event.Kv.Value)

				case clientv3.EventTypeDelete:
					a = namespace.WatcherDelete

					key = x.BytesToString(event.PrevKv.Key)
					b.Write(event.PrevKv.Value)

				}

				c <- &Action{Act: a, B: map[string][]byte{key: b.Bytes()}}
			}
		}
	}()

	return c
}

func (w *Watch) Close() {
	w.exit <- struct{}{}
	close(w.exit)
}
