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
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"roc/internal/etcd"
	"roc/internal/x"
	"roc/rlog"
	"strings"
	"sync"
)

//Configuration Center
//use etcd,

var gRConfig *config

var _ RConfig = &config{}

func init() {
	gRConfig = NewConfig()
}

type RConfig interface {

	// ConfigListAndSync Get all config and sync to cache
	ConfigListAndSync() error

	// WithConfig Get config with key
	WithConfig(key string) ([]byte, error)

	// SetConfig set config with key value
	SetConfig(key string, value []byte) error

	// Clean clean all config
	Clean() error

	// Delete remove a config with key
	Delete(key string) error

	// Watch watch config and update
	Watch() chan *etcd.Action

	// Backup backup config
	Backup() error

	// LoadFs2Etcd load a config file to etcd
	LoadFs2Etcd() error

	// Close close config
	Close()
}

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
}

func NewConfig(opts ...Options) *config {
	c := &config{
		opts:  newOpts(),
		data:  make(map[string][]byte),
		close: make(chan struct{}),
	}

	c.watch = etcd.NewEtcdWatch(c.opts.schema, c.opts.e.Client())
	c.action = c.watch.Watch(c.opts.schema)

	go c.run()

	return c
}

func (c *config) ConfigListAndSync() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	data, err := c.opts.e.GetWithList(c.opts.schema)
	if err != nil {
		return err
	}

	rlog.Infof("update config |data=%v", x.MustMarshalString(data))

	for k, v := range data {
		data[k] = v
	}

	retrun
	c.Backup()
}

func (c *config) WithConfig(key string) ([]byte, error) {
	return c.opts.e.GetWithKey(c.opts.schema + "/" + key)
}

func (c *config) SetConfig(key string, value []byte) error {
	return c.opts.e.Put(key, string(value))
}

// Clean clean all config,avoid using Clean if it is not necessary
func (c *config) Clean() error {
	return c.opts.e.Delete(c.opts.schema)
}

func (c *config) Delete(key string) error {
	return c.opts.e.Delete(c.opts.schema + "/" + key)
}

func (c *config) Backup() error {
	for k, v := range c.data {
		if name := getFsName(k); name != "" {
			//open or create a local config file
			fs, err := os.OpenFile(
				c.opts.backupPath+string(os.PathSeparator)+name,
				os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC,
				os.ModePerm)
			if err != nil {
				rlog.Error(err)
				continue
			}

			_, _ = fs.Write(v)

			_ = fs.Close()
		}
	}

	return nil
}

// LoadFs2Etcd load local config file to etcd
func (c *config) LoadFs2Etcd() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	dir, err := ioutil.ReadDir(c.opts.backupPath)
	if err != nil {
		return err
	}

	for _, v := range dir {
		if v.IsDir() {
			continue
		}

		b, err := ioutil.ReadFile(c.opts.backupPath + v.Name())
		if err != nil {
			rlog.Error(err)
			continue
		}

		c.data[c.opts.schema+v.Name()] = b

		err = c.SetConfig(c.opts.schema+v.Name(), b)
		if err != nil {
			rlog.Error(err)
		}
	}

	return nil
}

func getFsName(s string) string {
	array := strings.Split(s, "/")
	if len(array) > 0 {
		return array[len(array)-1]
	}

	return ""
}

// Watch watching and update all config
func (c *config) Watch() chan *etcd.Action {
	var r = make(chan *etcd.Action)
	go func() {
		for action := range c.action {
			r <- action
		}

		close(r)
	}()

	return r
}

func (c *config) run() {
	if c.opts.enableDynamic {
		for {
			select {
			case <-c.Watch():
				// sync config all
				c.lock.Lock()
				err := c.ConfigListAndSync()
				if err != nil {
					rlog.Error(err)
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

func RocConfig() *config {
	return gRConfig
}

func getDataBytes(args ...string) ([]byte, string) {
	var name, key string
	if len(args) == 1 {
		key = args[0]
		name = gRConfig.opts.schema + gRConfig.opts.prefix
	}

	if len(args) == 2 {
		key = args[1]
		name = args[0]
	}

	return gRConfig.data[name], key
}

func GetAny(args ...string) jsoniter.Any {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetAny no data error")
	}
	return x.Jsoniter.Get(b, key)
}

func GetString(args ...string) string {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		return ""
	}
	return x.Jsoniter.Get(b, key).ToString()
}

func GetInt(args ...string) int {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetInt no data error")
	}
	return x.Jsoniter.Get(b, key).ToInt()
}

func GetInt32(args ...string) int32 {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetInt no data error")
	}
	return x.Jsoniter.Get(b, key).ToInt32()
}

func GetInt64(args ...string) int64 {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetInt no data error")
	}
	return x.Jsoniter.Get(b, key).ToInt64()
}

func GetUint(args ...string) uint {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetInt no data error")
	}
	return x.Jsoniter.Get(b, key).ToUint()
}

func GetUint32(args ...string) uint32 {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetInt no data error")
	}
	return x.Jsoniter.Get(b, key).ToUint32()
}

func GetUint64(args ...string) uint64 {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetInt no data error")
	}
	return x.Jsoniter.Get(b, key).ToUint64()
}

func GetFloat32(args ...string) float32 {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetFloat32 no data error")
	}
	return x.Jsoniter.Get(b, key).ToFloat32()
}

func GetFloat64(args ...string) float64 {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetFloat64 no data error")
	}
	return x.Jsoniter.Get(b, key).ToFloat64()
}

func GetBool(args ...string) bool {
	b, key := getDataBytes(args...)
	if len(b) == 0 {
		panic("config GetBool no data error")
	}
	return x.Jsoniter.Get(b, key).ToBool()
}
