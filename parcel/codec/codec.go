package codec

import (
	"github.com/gogo/protobuf/proto"

	"roc/parcel/codec/protoc"
)

var DefaultCodec Codec = &protoc.Proto{}

type Codec interface {
	Encode(message proto.Message) ([]byte, error)
	Decode(b []byte, message proto.Message) error
}
