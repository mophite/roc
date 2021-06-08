package protoc

import (
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

func (*Proto) Decode(b []byte, rsp proto.Message) error {
	return proto.Unmarshal(b, rsp)
}
