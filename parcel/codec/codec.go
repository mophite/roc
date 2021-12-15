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

package codec

import (
    "github.com/gogo/protobuf/proto"
    "github.com/rsocket/rsocket-go/extension"

    "github.com/go-roc/roc/parcel/codec/jsonc"
    "github.com/go-roc/roc/parcel/codec/protoc"
)

//todo https://github.com/klauspost/compress use compress

var defaultCodec Codec = jsonc.JSCodec

type Codec interface {
    Encode(message proto.Message) ([]byte, error)
    Decode(b []byte, message proto.Message) error
    MustEncodeString(message proto.Message) string
    MustEncode(message proto.Message) []byte
    MustDecode(b []byte, message proto.Message)
    Name() string
}

var DefaultCodecs = map[string]Codec{
    extension.ApplicationJSON.String():     jsonc.JSCodec,
    extension.ApplicationProtobuf.String(): &protoc.Proto{},
}

func CodecType(contentType string) Codec {
    c, ok := DefaultCodecs[contentType]
    if !ok {
        return defaultCodec
    }
    return c
}

func SetCodec(contentType string, c Codec) {
    DefaultCodecs[contentType] = c
}
