package client

import (
	"roc/internal/namespace"
	"time"

	"github.com/gogo/protobuf/proto"

	"roc/internal/backoff"
	"roc/parcel"
	"roc/parcel/context"
)

type Invoke struct {
	strategy Strategy
	opts     InvokeOption
}

func newInvoke(c *context.Context, method string, client *RocClient, opts ...InvokeOptions) (*Invoke, error) {
	invoke := &Invoke{strategy: client.strategy}

	for i := range opts {
		opts[i](&invoke.opts)
	}

	if invoke.opts.buffSize == 0 {
		invoke.opts.buffSize = 10
	}

	var err = c.WithMetadata(
		invoke.opts.service,
		method,
		"",
		map[string]string{
			namespace.DefaultHeaderVersion: invoke.opts.version,
			namespace.DefaultHeaderAddress: invoke.opts.address,
		})
	return invoke, err
}

func (i *Invoke) invokeRR(c *context.Context, req, rsp proto.Message, conn *Conn, opts option) error {
	b, err := opts.cc.Encode(req)
	if err != nil {
		return err
	}

	var request, response = parcel.Payload(b), parcel.NewPacket()
	defer func() {
		parcel.Recycle(response, request)
		c.Info(req.String(), rsp.String())
	}()

	err = conn.Client().RR(c, request, response)
	if err != nil {
		bf := backoff.NewBackoff()
		for i := 0; i < opts.retry; i++ {
			time.Sleep(bf.Next(i))
			if err = conn.Client().RR(c, request, response); err == nil {
				break
			}
		}

		if err != nil {
			c.Error(err)
			conn.growError()
			return err
		}
	}

	return opts.cc.Decode(response.Bytes(), rsp)
}
