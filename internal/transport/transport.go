package transport

import (
	"roc/internal/endpoint"
	"roc/internal/router"
	"roc/parcel"
	"roc/parcel/context"
)

type Server interface {
	Address() string
	Accept(fn *router.Router)
	Run()
	String() string
	Close()
}

type Client interface {
	Dial(e *endpoint.Endpoint, closeChan chan string) error

	// RR request/response,through block unsafe method
	RR(c *context.Context, req *parcel.RocPacket, rsp *parcel.RocPacket) (err error)

	// FF FireAndForget
	//FF(c *context.Context, req *parcel.RocPacket)

	// RS request/stream
	RS(c *context.Context, req *parcel.RocPacket) (chan []byte, chan error)

	// RC request/channel
	RC(c *context.Context, req chan []byte, errsIn chan error) (chan []byte, chan error)

	// MP metadata
	//MP(c *context.Context)

	String() string

	Close()
}

type CallOptions func(option *CallOption)

type CallOption struct {
}
