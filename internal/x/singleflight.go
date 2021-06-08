package x

import "sync"

var group = &Group{m: make(map[string]*call)}

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
