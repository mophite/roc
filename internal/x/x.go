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

package x

import (
	"math/rand"
	"time"
	"unsafe"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func NewUUID() string {
	return uuid.New().String()
}

var Jsoniter = jsoniter.ConfigFastest

func MustMarshal(v interface{}) []byte {
	b, _ := Jsoniter.Marshal(v)
	return b
}

func MustUnmarshal2Map(b []byte, v map[string]string) {
	err := Jsoniter.Unmarshal(b, v)
	if err != nil {
		v = make(map[string]string)
	}
}

func MustMarshalString(v interface{}) string {
	b, _ := Jsoniter.MarshalToString(v)
	return b
}

func StringToBytes(s string) (b []byte) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return Rand.Intn(max-min) + min
}
