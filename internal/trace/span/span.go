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

package span

import (
    "github.com/go-roc/roc/x"
)

// attention:it's maybe zipkin instead,don't user this for your production Environment.

// Span
// |---TraceId:1     ----->RPC----->       |---TraceId:1
//       |---ParentSpanId:0                            |---ParentSpanId:222
//           |---SpanId:222                                 |---SpanId:223

// Span this is a demo,need to be richer
type Span struct {
    SpanId       uint32
    ParentSpanId int32
    traceId      string
}

func (s *Span) Carrier() {
    s.ParentSpanId += 1
    s.SpanId += 1
}

func (s *Span) Finish() {
    // todo buffer flush to cloud or something
    return
}

func (s *Span) Name() string {
    return "span"
}

func (s *Span) TraceId() string {
    return s.traceId
}

func NewSpan() *Span {
    return &Span{
        traceId:      x.NewXID(),
        ParentSpanId: -1,
        SpanId:       1,
    }
}
