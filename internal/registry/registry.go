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

package registry

import (
    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/namespace"
)

type Registry interface {

    // Watcher watch the remote like etcd,who's data change
    Watcher

    // Register register a endpoint to etcd or other
    Register(e *endpoint.Endpoint) error

    // Deregister deregister a endpoint from etcd or other
    Deregister(e *endpoint.Endpoint) error

    // Next get one endpoint
    Next(scope string) (*endpoint.Endpoint, error)

    // List get all endpoint
    List() ([]*endpoint.Endpoint, error)

    // Name return the tool's name like "etcd"
    Name() string

    // CloseRegistry close registry
    CloseRegistry()
}

type Watcher interface {
    Watch() chan *Action
}

// Action watch data change content
type Action struct {

    Act namespace.WatcherAction

    E   *endpoint.Endpoint

    //etcd key
    Key string
}
