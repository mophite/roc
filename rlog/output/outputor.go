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

package output

import (
    "bytes"

    "github.com/go-roc/roc/rlog/common"
    "github.com/go-roc/roc/rlog/output/console"
)

type Outputor interface {
    Init(string)
    Out(level common.Level, b *bytes.Buffer)
    Level() common.Level
    SetLevel(level common.Level)
    Poller()
    Close()
    String() string
}

var DefaultOutput Outputor = &console.Console{L: common.DEBUG}
