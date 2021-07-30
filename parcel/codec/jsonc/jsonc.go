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

package jsonc

import (
    "github.com/go-roc/roc/x"
    "github.com/gogo/protobuf/proto"
    jsoniter "github.com/json-iterator/go"
)

var JSCodec = &JsonCodec{jsonCodec: x.Jsoniter}

//var JSCodec = &JsonCodec{jsonCodec: jsoniter.ConfigCompatibleWithStandardLibrary}

type JsonCodec struct {
    jsonCodec jsoniter.API
}

func (j *JsonCodec) Encode(req proto.Message) ([]byte, error) {
    b, err := j.jsonCodec.Marshal(req)
    if err != nil {
        return nil, err
    }

    return b, nil
}

func (j *JsonCodec) MustEncode(req proto.Message) []byte {
    b, _ := j.jsonCodec.Marshal(req)
    return b
}

func (j *JsonCodec) MustEncodeString(req proto.Message) string {
    b, err := j.jsonCodec.MarshalToString(req)
    if err != nil {
        return ""
    }

    return b
}

func (j *JsonCodec) Decode(b []byte, rsp proto.Message) error {
    return j.jsonCodec.Unmarshal(b, rsp)
}

func (j *JsonCodec) MustDecode(b []byte, rsp proto.Message) {
    _ = j.jsonCodec.Unmarshal(b, rsp)
}

func (j *JsonCodec) Name() string {
    return "jsoniter"
}
