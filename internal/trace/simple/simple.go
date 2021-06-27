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

package simple

import (
    "github.com/go-roc/roc/x"
)

// Simple |---TraceId:1     ----->RPC----->       |---TraceId:1
type Simple struct {
    traceId string
}

func (s *Simple) Carrier() {
    return
}

func NewSimple(traceId ...string) *Simple {
    s := &Simple{}
    if len(traceId) > 0 && traceId[0] != "" {
        s.traceId = traceId[0]
        return s
    }
    s.traceId = x.NewUUID()
    return s
}

func (s *Simple) String() string {
    return "simple"
}

func (s *Simple) Finish() {
    return
}

func (s *Simple) TraceId() string {
    return s.traceId
}
