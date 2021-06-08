package im

import (
	"fmt"
	"roc/_auxiliary/example/tutorials/proto/pbim"
	"roc/parcel/context"
)

type Im struct {
	H *Hub
	p *point
}

func (i *Im) Connect(c *context.Context, req *pbim.ConnectReq) (rsp *pbim.ConnectRsp, err error) {
	i.p = &point{userName: req.UserName, message: make(chan *pbim.SendMessageRsp)}
	i.H.addClient(i.p)
	return &pbim.ConnectRsp{IsConnect: true}, nil
}

func (i *Im) Count(c *context.Context, req *pbim.CountReq) (rsp *pbim.CountRsp, err error) {
	return &pbim.CountRsp{
		Count: i.H.count(),
	}, nil
}

func (i *Im) SendMessage(c *context.Context, req chan *pbim.SendMessageReq, errIn chan error) (chan *pbim.SendMessageRsp, chan error) {
	var rsp = make(chan *pbim.SendMessageRsp)
	var errs = make(chan error)

	go func() {
		for data := range i.p.message {
			rsp <- data
		}
		close(rsp)
	}()

	go func() {
	QUIT:
		for {
			select {
			case data, ok := <-req:
				if !ok {
					break QUIT
				}

				i.H.broadCast <- data

			case e := <-errIn:
				if e != nil {
					fmt.Println("----------close------")
					errs <- e
					break QUIT
				}
			}
		}

		close(errs)
		i.H.removeClient(i.p)
	}()

	return rsp, errs
}
