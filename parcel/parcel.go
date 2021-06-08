package parcel

import (
	"github.com/gogo/protobuf/proto"

	"roc/parcel/context"
)

type Handler func(c *context.Context, req *RocPacket, interrupt Interceptor) (rsp proto.Message, err error)

type StreamHandler func(c *context.Context, req *RocPacket) (chan proto.Message, chan error)

type ChannelHandler func(c *context.Context, req chan *RocPacket, errs chan error) (chan proto.Message, chan error)

type Fire func(c *context.Context, req proto.Message) (proto.Message, error)

type Interceptor func(c *context.Context, req proto.Message, fire Fire) (proto.Message, error)

type Wrapper func(c *context.Context) error
