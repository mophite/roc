package main

import (
	"context"
	"log"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)

func main() {
	tp := rsocket.TCPServer().SetHostAndPort("127.0.0.1", 10000).Build()
	err := rsocket.Receive().
		OnStart(func() {
			log.Println("server start success! at:", 10000)
		}).
		Acceptor(func(
			ctx context.Context,
			setup payload.SetupPayload,
			sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {

			sendingSocket.OnClose(func(err error) {
				log.Println("*** socket disconnected ***")
			})

			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(pl payload.Payload) mono.Mono {
					return mono.JustOneshot(pl)
				}),
			), nil
		}).
		Transport(tp).
		Serve(context.Background())
	if err != nil {
		panic(err)
	}
}
