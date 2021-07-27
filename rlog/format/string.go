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

    "github.com/go-roc/roc/rlog/common"
    "github.com/go-roc/roc/x/bytesbuffpool"
)

var _ Formatter = &stringFormat{}

type stringFormat struct {
	layout string
}

func (s *stringFormat) Layout() string {
	if s.layout == "" {
		return defaultLayout
	}
	return s.layout
}

func (s *stringFormat) Format(detail *common.Detail) *bytes.Buffer {
	b := bytesbuffpool.Get()

	b.WriteString("[" + detail.Level + "] ")

	if detail.Line != "" {
		b.WriteString(detail.Line + " ")
	}

	b.WriteString(detail.Timestamp + " ")

	b.WriteString(detail.Content)

	return b
}

func (s *stringFormat) SetLayout(layout string) {
	s.layout = layout
}

func (s *stringFormat) String() string {
	return "str"
}
