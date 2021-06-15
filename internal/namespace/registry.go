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

package namespace

import (
	"strings"
)

// SplicingPrefix is splicing service address
// Schema default is DefaultSchema. eg. goroc
func SplicingPrefix(schema Schema, scope Scope) string {
	var b strings.Builder
	b.WriteString(schema)

	if scope != "" {
		b.WriteString("/")
		b.WriteString(scope)
	}

	return b.String()
}

// SplicingScope is the service name/version. eg.srv.hello/version
// name or version don't allowed none
func SplicingScope(name, version string) Scope {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString("/")
	b.WriteString(version)
	return b.String()
}
