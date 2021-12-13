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
)

const defaultLayout = "2006.01.02 15:04:05.000"

var DefaultFormat Formatter = &stringFormat{}

type Formatter interface {
	Layout() string
	Format(detail *common.Detail) *bytes.Buffer
	SetLayout(layout string)
	String() string
}
