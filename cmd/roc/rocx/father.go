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

package rocx

import (
	"os"
	"path/filepath"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

// getLastPwd get last file directory
func getLastPwd() string {
	f := func(s string, pos, length int) string {
		runes := []rune(s)
		l := pos + length
		if l > len(runes) {
			l = len(runes)
		}
		return string(runes[pos:l])
	}
	directory := pwd()
	return f(directory, 0, strings.LastIndex(directory, pathSeparator))
}

// pwd get current file directory
func pwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

// getProjectName get current project name
func getProjectName() string {
	f := func(s string, pos int) string {
		runes := []rune(s)
		return string(runes[pos:])
	}

	directory := pwd()
	return strings.Trim(f(directory, strings.LastIndex(directory, pathSeparator)), pathSeparator)
}
