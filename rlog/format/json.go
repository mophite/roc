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

package format

import (
    "bytes"

    "github.com/go-roc/roc/x/bytesbuffpool"
    jsoniter "github.com/json-iterator/go"

    "github.com/go-roc/roc/rlog/common"
)

var _ Formatter = &jsonFormat{}

type jsonFormat struct {
	layout string
}

func (j *jsonFormat) Layout() string {
	if j.layout == "" {
		return defaultLayout
	}
	return j.layout
}

func (j *jsonFormat) Format(detail *common.Detail) *bytes.Buffer {
	b := bytesbuffpool.Get()
	b.Write(mustMarshal(detail))
	return b
}

func (j *jsonFormat) SetLayout(layout string) {
	j.layout = layout
}

func (j *jsonFormat) String() string {
	return "json"
}

var fastest = jsoniter.ConfigFastest

func mustMarshal(v interface{}) []byte {
	b, _ := fastest.Marshal(v)
	return b
}
