package bench_rpc

import (
    "fmt"
    "sync/atomic"
    "time"

    "github.com/go-roc/roc/_auxiliary/example/tutorials/internal/ipc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

func RPC() error {
    _, err := ipc.SaySrv(context.Background(), &phello.SayReq{})
    return err
}

func BenchRpc() {

    var (
        tps      int64
        errCount int64
    )
    for i := 0; i < 500; i++ {
        go func() {
            for j := 0; j < 100000; j++ {
                err := RPC()
                if err != nil {
                    atomic.AddInt64(&errCount, 1)
                    continue
                }
                atomic.AddInt64(&tps, 1)
            }
        }()
    }

    for i := 0; i < 10; i++ {
        time.Sleep(time.Second)
        fmt.Printf(
            "----------------------tps=%v err=%v------------------\n",
            atomic.LoadInt64(&tps),
            atomic.LoadInt64(&errCount),
        )
        atomic.SwapInt64(&tps, 0)
        atomic.SwapInt64(&errCount, 0)
    }
}

//mac
//----------------------tps=68673 err=0------------------
//----------------------tps=71198 err=0------------------
//----------------------tps=74811 err=0------------------
//----------------------tps=71500 err=0------------------
//----------------------tps=73527 err=0------------------
//----------------------tps=70633 err=0------------------
//----------------------tps=69648 err=0------------------
//----------------------tps=74207 err=0------------------
//----------------------tps=68339 err=0------------------
//----------------------tps=69578 err=0------------------
