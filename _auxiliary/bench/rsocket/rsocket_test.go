package main

import (
    "bytes"
    "context"
    "fmt"
    "runtime"
    "sync"
    "sync/atomic"
    "testing"
    "time"

    "github.com/jjeffcaii/reactor-go/scheduler"
    "github.com/rsocket/rsocket-go"
    "github.com/rsocket/rsocket-go/payload"
    "github.com/rsocket/rsocket-go/rx/mono"
)

var cc = make(chan struct{})

var c rsocket.Client

func init() {
    go func() {
        err := RSocketServer()
        if err != nil {
            panic(err)
        }
    }()

    <-cc

    var err error
    c, err = createClient()
    if err != nil {
        panic(err)
    }
}

var count int64

// RSocketServer is a simple rsocket server.
func RSocketServer() error {

    return rsocket.Receive().
        Resume().
        Fragment(4096).OnStart(
        func() {
            cc <- struct{}{}
        },
    ).
        Acceptor(
            func(
                ctx context.Context,
                setup payload.SetupPayload,
                sendingSocket rsocket.CloseableRSocket,
            ) (rsocket.RSocket, error) {
                r := rsocket.NewAbstractSocket(
                    rsocket.RequestResponse(
                        func(msg payload.Payload) mono.Mono {
                            return mono.Just(msg)
                        },
                    ),
                    rsocket.FireAndForget(
                        func(request payload.Payload) {
                            atomic.AddInt64(&count, 1)
                        },
                    ),
                )
                return r, nil
            },
        ).
        Transport(rsocket.TCPServer().SetHostAndPort("127.0.0.1", 11111).Build()).
        Serve(context.Background())
}

func createClient() (rsocket.Client, error) {
    return rsocket.Connect().
        Scheduler(scheduler.NewElastic(runtime.NumCPU()*20), nil).
        //Lease().
        //Resume().
        KeepAlive(time.Second*1, time.Second*10, 1).
        ConnectTimeout(time.Second).
        OnClose(
            func(err error) {
                fmt.Println("*** disconnected ***")
            },
        ).
        SetupPayload(payload.NewString("你好", "世界")).
        Acceptor(
            func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
                return rsocket.NewAbstractSocket(
                    rsocket.RequestResponse(
                        func(p payload.Payload) mono.Mono {
                            fmt.Println("receive request from server:", p)
                            if bytes.Equal(p.Data(), []byte("ping")) {
                                return mono.Just(payload.NewString("pong", "from client"))
                            }
                            return mono.Just(p)
                        },
                    ),
                )
            },
        ).
        Transport(rsocket.TCPClient().SetHostAndPort("127.0.0.1", 11111).Build()).
        Start(context.Background())
}

func BenchmarkRsocketServer(b *testing.B) {

    b.ResetTimer()

    var wg sync.WaitGroup
    for i := 0; i < b.N; i++ {
        wg.Add(1)
        go func() {
            _, release, err := c.RequestResponse(payload.NewString("1", "")).BlockUnsafe(context.TODO())
            if err != nil {
                panic(err)
            }
            release()
            wg.Done()
        }()
    }

    wg.Wait()
}

func TestRsocketRR(t *testing.T) {

    var tps int64
    for i := 0; i < 100; i++ {
        go func() {
            for j := 0; j < 100000; j++ {
                _, err := c.RequestResponse(payload.NewString("1", "")).Block(context.TODO())
                if err != nil {
                    panic(err)
                }
                atomic.AddInt64(&tps, 1)
            }
        }()
    }

    var tmp int64
    for i := 0; i < 10; i++ {
        time.Sleep(time.Second)
        t1 := atomic.LoadInt64(&tps)
        fmt.Println("----------tps--------", t1)
        tmp += t1
        atomic.SwapInt64(&tps, 0)
    }
}

func BenchmarkRsocketServerNoRsp(b *testing.B) {

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        c.FireAndForget(payload.NewString("1", ""))
    }
    fmt.Println("--------------", atomic.LoadInt64(&count))
}

func TestRsocketFF(t *testing.T) {

    var tps int64
    for i := 0; i < 100; i++ {
        go func() {
            for j := 0; j < 10000; j++ {
                c.FireAndForget(payload.NewString("1", ""))
                atomic.AddInt64(&tps, 1)
            }
        }()
    }

    var tmp int64
    for i := 0; i < 3; i++ {
        time.Sleep(time.Second)
        t1 := atomic.LoadInt64(&tps)
        fmt.Println("----------tps--------", t1)
        tmp += t1
        atomic.SwapInt64(&tps, 0)
    }

    time.Sleep(time.Second * 5)

    fmt.Println("-------rsp-------", atomic.LoadInt64(&count))
    fmt.Println("-------actual-------", tmp)
}
