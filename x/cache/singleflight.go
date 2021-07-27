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

package cache

import "sync"

var group = &Group{m: make(map[string]*call)}

//user for cache coherency
//When the cache expires at a certain point in time,
//there are a large number of concurrent requests for this key at this point in time.
//These requests find that the cache expires generally will load data from the backend DB and reset it to the cache.
//At this time,
//there are large concurrent requests. It may overwhelm the back-end DB in an instant.
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	group.mu.Lock()
	if group.m == nil {
		group.m = make(map[string]*call)
	}
	if c, ok := group.m[key]; ok {
		group.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	group.m[key] = c
	group.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	group.mu.Lock()
	delete(group.m, key)
	group.mu.Unlock()

	return c.val, c.err
}
