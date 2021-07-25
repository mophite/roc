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

package endpoint

import (
    "errors"

    "github.com/go-roc/roc/internal/namespace"
)

// DefaultLocalEndpoint service discovery endpoint
var DefaultLocalEndpoint *Endpoint

type Endpoint struct {
    //endpoint unique id
    Id string

    //endpoint name
    Name string

    //endpoint version
    Version string

    // schema/name/version/id
    Absolute string

    //service server ip address
    Address string

    // name.version
    // eg. api.hello/v.1.0.0
    Scope string
}

// NewEndpoint new a endpoint with schema,id,name,version,address
func NewEndpoint(id, name, address string) (*Endpoint, error) {
    if name == "" || address == "" || id == "" {
        return nil, errors.New("not complete")
    }
    e := new(Endpoint)
    e.Id = id
    e.Name = name
    e.Version = namespace.DefaultVersion
    e.Address = address
    e.Scope = e.Name + "/" + e.Version
    e.Absolute = namespace.DefaultSchema + "/" + e.Scope + "/" + e.Address
    return e, nil
}
