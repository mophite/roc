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

package fs

import (
	"errors"
	"os"
	"path"
	"strings"
)

func open(name string, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, perm)
}

func isLink(filename string) (string, error) {
	fi, err := os.Lstat(filename)
	if err != nil {
		return "", err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		name, err := os.Readlink(filename)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	return "", errors.New("not symlink")
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func getFilenamePrefix(s string) string {
	return strings.TrimSuffix(path.Base(s), ".log")
}
