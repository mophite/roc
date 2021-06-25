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
    "bytes"
    "context"
    "strings"

    "github.com/coreos/etcd/clientv3"

    "github.com/go-roc/roc/internal/x"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/rlog"
)

type Action struct {

    //watch callback action
    Act namespace.WatcherAction

    //callback data
    B map[string][]byte
}

type Watch struct {

    //exit signal
    exit chan struct{}

    //etcd client
    client *clientv3.Client

    //etcd watch channel
    wc clientv3.WatchChan
}

func NewEtcdWatch(prefix string, client *clientv3.Client) *Watch {
    var w = &Watch{
        exit:   make(chan struct{}),
        client: client,
    }

    w.wc = client.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
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

            //watch result events
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

// Close close watch
func (w *Watch) Close() {
    w.exit <- struct{}{}
    close(w.exit)
}
