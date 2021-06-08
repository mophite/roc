package parcel

import (
	"roc/parcel/codec"

	"roc/parcel/packet"
)

type ErrorCode = int32

const (
	ErrorCodeInternalServer ErrorCode = 500
	ErrorCodeBadRequest               = 400
	ErrCodeNotFoundHandler            = 401
)

var (
	DefaultErrorPacket = NewErrorPacket(codec.DefaultCodec)
)

type ErrorPackager interface {
	Encode(code ErrorCode, err error) []byte
}

type ErrorPacket struct {
	data packet.Packet
	cc   codec.Codec
}

func NewErrorPacket(cc codec.Codec) *ErrorPacket {
	return &ErrorPacket{cc: cc}
}

func (e *ErrorPacket) Encode(code ErrorCode, err error) []byte {
	ep := *e
	ep.data.Code = code
	ep.data.Message = err.Error()
	b, _ := e.cc.Encode(&ep.data)
	return b
}
