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

package context

import (
	"github.com/go-roc/roc/internal/trace/simple"
	"github.com/go-roc/roc/parcel/metadata"
)

const (
	//default packet pool size
	defaultPoolSize = 10240000
)

//create default pool
var pool = &contextPool{
	c: make(chan *Context, defaultPoolSize),
}

type contextPool struct {
	c chan *Context
}

func Recycle(p *Context) {

	p.reset()

	select {
	case pool.c <- p:
	default: //if pool full,throw away
	}
}

func New() (p *Context) {
	select {
	case p = <-pool.c:
		p.Trace = simple.NewSimple()
		p.Metadata = metadata.MallocMetadata()
		p.data = make(map[string]interface{}, 10)
	default:
		p = newContext()
	}

	return
}
