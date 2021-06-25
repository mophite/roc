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
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"

	"github.com/go-roc/roc/etcd"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/x"
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

	cache map[string]interface{}
}

func NewConfig(opts ...Options) error {
	gRConfig = &config{
		opts:  newOpts(),
		data:  make(map[string][]byte),
		cache: make(map[string]interface{}),
		close: make(chan struct{}),
	}

	gRConfig.watch = etcd.NewEtcdWatch(gRConfig.opts.schema, gRConfig.opts.e.Client())
	gRConfig.action = gRConfig.watch.Watch(gRConfig.opts.schema)

	err := gRConfig.configListAndSync()
	if err != nil {
		rlog.Error(err)
		return err
	}

	go gRConfig.update()

	return nil
}

func (c *config) configListAndSync() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	globalData, err := c.opts.e.GetWithList(c.opts.global, clientv3.WithPrefix())
	if err == nil {
		for k, v := range globalData {
			c.data[getFsName(k)] = v
		}
	}

	privateData, err := c.opts.e.GetWithList(c.opts.private, clientv3.WithPrefix())
	if err == nil {
		for k, v := range privateData {

			//cover global config
			if _, ok := c.data[getFsName(k)]; ok {
				c.data[getFsName(k)] = v
				continue
			}

			c.data[getFsName(k)] = v
		}
	}

	return c.backup()
}

func (c *config) backup() error {
	fs, err := os.OpenFile(
		c.opts.backupPath,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.ModePerm,
	)

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

	fs.Write(b)

	fs.Close()

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

	return x.Jsoniter.Unmarshal(fd, &c.data)
}

func getFsName(s string) string {
	isGlobal := strings.Contains(s, gRConfig.opts.global)
	array := strings.Split(s, "/")

	if len(array) > 0 {
		s = array[len(array)-1]
	}

	if isGlobal {
		s = gRConfig.opts.prefix + "." + s
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
				case namespace.WatcherCreate:
					for k, v := range data.B {

						var key = getFsName(k)

						if _, ok := c.data[key]; !ok {
							c.data[key] = v
							if f, ok := c.cache[key]; ok {
								err := Decode2Config(key, f)
								if err != nil {
									rlog.Error(err)
								}
							}
						} else {
							rlog.Warnf("same config warning: %s", key)
						}
					}

				case namespace.WatcherUpdate:
					for k, v := range data.B {

						var key = getFsName(k)

						if _, ok := c.data[key]; ok {
							c.data[key] = v
							if f, ok := c.cache[key]; ok {
								err := Decode2Config(key, f)
								if err != nil {
									rlog.Error(err)
								}
							}
						} else {
							rlog.Warnf("config not exist: %s", key)
						}
					}

				case namespace.WatcherDelete:
					for k, _ := range data.B {

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

func (c *config) Close() {
	c.lock.Lock()
	c.data = nil
	c.lock.Unlock()
	c.close <- struct{}{}
}

func getDataBytes(key string) []byte {
	return gRConfig.data[key]
}

func Decode2Config(key string, v interface{}) error {
	err := x.Jsoniter.Unmarshal(getDataBytes(key), v)
	if err != nil {
		return err
	}

	gRConfig.cache[key] = v

	return nil
}
