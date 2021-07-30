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

package protoc

import (
    "github.com/go-roc/roc/x"
    "github.com/gogo/protobuf/proto"
)

type Proto struct{}

func (p *Proto) Encode(req proto.Message) ([]byte, error) {
    b, err := proto.Marshal(req)
    if err != nil {
        return nil, err
    }

    return b, nil
}

func (p *Proto) MustEncodeString(req proto.Message) string {
    b, err := proto.Marshal(req)
    if err != nil {
        return ""
    }
    return x.BytesToString(b)
}

func (p *Proto) MustEncode(req proto.Message) []byte {
    b, _ := proto.Marshal(req)
    return b
}

func (*Proto) Decode(b []byte, rsp proto.Message) error {
    err := proto.Unmarshal(b, rsp)
    if err != nil {
        return err
    }
    return nil
}

func (*Proto) MustDecode(b []byte, rsp proto.Message) {
    proto.Unmarshal(b, rsp)
}

func (*Proto) Name() string {
    return "gogo_proto"
}
