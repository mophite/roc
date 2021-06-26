package main

import (
	"context"
	"fmt"

	"github.com/rsocket/rsocket-go"
)

func main() {
	createClient(20000)
	createClient(10000)

	select {}
}

func createClient(port int) (rsocket.Client, error) {
	return rsocket.Connect().
		OnClose(
			func(err error) { //when net occur some error,it's will be callback the error server ip address
				fmt.Println("-----------关闭的连接-------", port)
			},
		).
		OnConnect(func(client rsocket.Client, err error) {
			fmt.Println("---------连接", port)
		}).
		Transport(rsocket.TCPClient().SetHostAndPort("127.0.0.1", port).Build()).
		Start(context.Background())
}
