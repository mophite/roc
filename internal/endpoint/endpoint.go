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
	"roc/internal/namespace"
)

type Endpoint struct {
	//endpoint unique id
	Id string

	//endpoint name
	Name string

	//enpoint version
	Version string

	// schema.name.version.id
	// eg. goroc.api.hello.v.1.1.1.2d1bd2f9-6951-4235-83bd-d6f38b358552
	Absolute string

	//service server ip address
	Address string

	// name.version
	// eg. api.hello.v.1.1.1
	Scope string
}

// Splicing generate scope and absolute with schema
func (e *Endpoint) Splicing(schema string) *Endpoint {
	e.Scope = namespace.SplicingScope(e.Name, e.Version)
	e.Absolute = schema + "/" + e.Scope + "/" + e.Address
	return e
}
