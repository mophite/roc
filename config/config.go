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

package config

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "sync"

    "github.com/coreos/etcd/clientv3"
    "github.com/go-roc/roc/x"

    "github.com/go-roc/roc/internal/etcd"
    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/rlog"
)

//Configuration Center
//use etcd,
var gRConfig *config

type config struct {

    //config option
    opts Option

    lock sync.RWMutex

    //config data local cache
    data map[string][]byte

    //close signal
    close chan struct{}

    //receive etcd callback data
    action chan *etcd.Action

    //watch etcd changed
    watch *etcd.Watch

    //cache the config objects
    //when create or update action occur
    //them will be updated
    cache map[string]interface{}
}

func NewConfig(opts ...Options) error {
    gRConfig = &config{
        opts:  newOpts(opts...),
        data:  make(map[string][]byte),
        cache: make(map[string]interface{}),
        close: make(chan struct{}),
    }

    if !gRConfig.opts.disableDynamic {
        gRConfig.watch = etcd.NewEtcdWatch(gRConfig.opts.schema, gRConfig.opts.e.Client())
        gRConfig.action = gRConfig.watch.Watch(gRConfig.opts.schema)
    }

    if gRConfig.opts.localFile {
        err := gRConfig.loadLocalFile()
        if err != nil {
            rlog.Error(err)
        }
    } else {
        err := gRConfig.configListAndSync()
        if err != nil {
            rlog.Error(err)
            return err
        }
    }

    if !gRConfig.opts.disableDynamic {
        go gRConfig.update()
    }

    if len(gRConfig.opts.f) > 0 {
        for i := range gRConfig.opts.f {
            err := gRConfig.opts.f[i]()
            if err != nil {
                rlog.Error(err)
                return err
            }
        }
    }

    return nil
}

func (c *config) configListAndSync() error {
    c.lock.Lock()
    defer c.lock.Unlock()

    publicData, err := c.opts.e.GetWithList(c.opts.public, clientv3.WithPrefix())
    if err == nil {
        for k, v := range publicData {
            c.data[getFsName(k)] = v
        }
    }

    privateData, err := c.opts.e.GetWithList(c.opts.private, clientv3.WithPrefix())
    if err == nil {
        for k, v := range privateData {
            c.data[getFsName(k)] = v
        }
    }

    return c.backup()
}

// PutPublic put public key value to etcd and cache data
func PutPublic(key, value string) error {
    c := gRConfig
    c.lock.Lock()
    defer c.lock.Unlock()

    err := c.opts.e.Put(c.opts.public+c.opts.publicPrefix+key, value)
    if err != nil {
        return err
    }
    c.data[c.opts.publicPrefix+key] = []byte(value)

    return nil
}

// PutPrivate put private key value to etcd and cache data
func PutPrivate(key, value string) error {
    c := gRConfig
    c.lock.Lock()
    defer c.lock.Unlock()

    err := c.opts.e.Put(c.opts.private+key, value)
    if err != nil {
        return err
    }

    c.data[key] = []byte(value)

    return nil
}

func (c *config) backup() error {
    fs, err := os.OpenFile(c.opts.backupPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)

    if err != nil {
        return err
    }

    var data = make(map[string]interface{})
    for k, v := range c.data {
        var tmp = make(map[string]interface{})
        err = x.Jsoniter.Unmarshal(v, &tmp)
        if err != nil {
            continue
        }
        data[k] = tmp
    }

    b, err := x.Jsoniter.Marshal(data)
    if err != nil {
        return err
    }

    _, _ = fs.Write(b)

    _ = fs.Close()

    return nil
}

// loadLocalFile load local config file to etcd
func (c *config) loadLocalFile() error {
    c.lock.Lock()
    defer c.lock.Unlock()

    fs, err := os.Open(c.opts.backupPath)
    if err != nil {
        return err
    }
    fd, err := ioutil.ReadAll(fs)
    if err != nil {
        return err
    }

    if gRConfig.opts.logOut {
        rlog.Infof("loadLocalFile |data=%s", string(fd))
    }

    var tmp = make(map[string]interface{})
    err = x.Jsoniter.Unmarshal(fd, &tmp)
    if err != nil {
        return err
    }

    for i := range tmp {
        c.data[i] = x.MustMarshal(tmp[i])
    }

    return nil
}

func getFsName(s string) string {
    array := strings.Split(s, "/")

    if len(array) > 0 {
        s = array[len(array)-1]
    }

    return s
}

func (c *config) update() {
    if !c.opts.disableDynamic {
        for {
            select {
            case data := <-c.action:
                // sync config all
                c.lock.Lock()

                switch data.Act {
                case namespace.WatcherCreate, namespace.WatcherUpdate:
                    for k, v := range data.B {

                        var key = getFsName(k)

                        c.data[key] = v

                        //load create config or update config to exist object
                        if f, ok := c.cache[key]; ok { //if ok,load to object
                            var err error
                            if strings.Contains(key, c.opts.publicPrefix) {
                                err = DecodePublic(key, f)
                            } else {
                                err = DecodePrivate(key, f)
                            }

                            if err != nil {
                                rlog.Error(err)
                            }
                        }
                    }

                case namespace.WatcherDelete:
                    for k := range data.B {

                        var key = getFsName(k)
                        if _, ok := c.data[key]; ok {
                            delete(c.data, key)
                            delete(c.cache, key)
                        }
                    }
                }

                c.lock.Unlock()

            case <-c.close:
                return
            }
        }
    }
}

func Close() {
    gRConfig.lock.Lock()
    gRConfig.data = nil
    gRConfig.lock.Unlock()
    gRConfig.close <- struct{}{}
}

// DecodePublic decode data to config and config will be updated when etcd watch change.
func DecodePublic(key string, v interface{}) error {

    d, ok := gRConfig.data[gRConfig.opts.publicPrefix+key]
    if !ok {
        return fmt.Errorf("config: %s not found", key)
    }
    err := x.Jsoniter.Unmarshal(d, v)
    if err != nil {
        return err
    }

    gRConfig.cache[key] = v

    if gRConfig.opts.logOut {
        rlog.Infof("DecodePublic |key=%s |value=%s", key, x.MustMarshalString(v))
    }

    return nil
}

// DecodePrivate decode data to config and config will be updated when etcd watch change.
func DecodePrivate(key string, v interface{}) error {

    d, ok := gRConfig.data[key]
    if !ok {
        return fmt.Errorf("config: %s not found", key)
    }

    err := x.Jsoniter.Unmarshal(d, v)
    if err != nil {
        return err
    }

    gRConfig.cache[key] = v

    if gRConfig.opts.logOut {
        rlog.Infof("DecodePrivate |key=%s |value=%s", key, x.MustMarshalString(v))
    }

    return nil
}
