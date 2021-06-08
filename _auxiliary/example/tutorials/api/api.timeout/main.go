package main

import (
	"fmt"
	"roc/_auxiliary/example/tutorials/proto/pbhello"
	"roc/client"
	"roc/parcel/context"
	"sync/atomic"
	"time"
)

var helloClient = pbhello.NewHelloWorldClient(client.NewRocClient())
var opt = client.WithScope("srv.hello")

func main() {

	rsp, errs := helloClient.SayStream(context.Background(), &pbhello.SayReq{Inc: 1}, opt, client.Timeout(time.Second*1))

	var count uint32

	var done = make(chan struct{})
	go func() {
		var err error
	QUIT:
		for {
			select {
			case b, ok := <-rsp:
				if ok {
					fmt.Println("------receive from srv.hello----", b.Inc)
					atomic.AddUint32(&count, 1)
				} else {
					break QUIT
				}
			case err = <-errs:
				if err != nil {
					break QUIT
				}
			}
		}
		done <- struct{}{}

		fmt.Println("say handler count is: ", atomic.LoadUint32(&count))
	}()

	<-done

}
