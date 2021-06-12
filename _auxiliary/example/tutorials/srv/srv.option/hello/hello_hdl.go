package hello

import (
	"roc/_auxiliary/example/tutorials/proto/pbhello"
	"roc/client"
	"roc/parcel/context"
	"sync/atomic"
	"time"
)

type Hello struct {
	Client *client.RocClient
}

func (h *Hello) SayStream(c *context.Context, req *pbhello.SayReq) (chan *pbhello.SayRsp, chan error) {
	var rsp = make(chan *pbhello.SayRsp)
	var err = make(chan error)

	go func() {
		var count uint32
		for i := 0; i < 200; i++ {
			rsp <- &pbhello.SayRsp{Inc: req.Inc + uint32(i)}
			atomic.AddUint32(&count, 1)
			time.Sleep(time.Second * 1)
		}

		c.Info("say stream example count is: ", atomic.LoadUint32(&count))

		close(rsp)
		close(err)
	}()

	return rsp, err
}

func (h *Hello) SayChannel(c *context.Context, req chan *pbhello.SayReq, errIn chan error) (chan *pbhello.SayRsp, chan error) {
	var rsp = make(chan *pbhello.SayRsp)
	var errs = make(chan error)

	go func() {
	QUIT:
		for {
			select {
			case data, ok := <-req:
				if !ok {
					break QUIT
				}

				//test channel sending frequency
				time.Sleep(time.Second)
				rsp <- &pbhello.SayRsp{Inc: data.Inc + uint32(1)}

			case e := <-errIn:
				if e != nil {
					errs <- e
				}
			}
		}

		close(rsp)
		close(errs)
	}()

	return rsp, errs
}

func (h *Hello) Say(c *context.Context, req *pbhello.SayReq) (rsp *pbhello.SayRsp, err error) {
	//  when set timeout is time.Second*1,it's will occur cancelled error
	time.Sleep(time.Second * 2)
	return &pbhello.SayRsp{Inc: req.Inc + 1}, nil
}
