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

package console

import (
    "bytes"
    "fmt"

    "github.com/go-roc/roc/rlog/common"
)

type Console struct {
    L common.Level
}

func (s *Console) Init(string) {
    return
}

func (s *Console) Out(level common.Level, b *bytes.Buffer) {

    if level < s.L {
        return
    }

    fmt.Printf(b.String())

    common.Buffer.Put(b)
}

func (s *Console) Level() common.Level {
    return s.L
}

func (s *Console) SetLevel(l common.Level) {
    s.L = l
}

func (s *Console) Poller() {
    return
}

func (s *Console) Close() {
    return
}

func (s *Console) String() string {
    return "console"
}
