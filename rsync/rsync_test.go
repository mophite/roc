package rsync

import (
    "fmt"
    "testing"
    "time"

    _ "github.com/go-roc/roc/etcd/mock"
)

func TestAcquire(t *testing.T) {
    const key = "test"
    go func() {
        err := Acquire(
            key, 30, 3, func() error {
                fmt.Println("do something!")
                time.Sleep(time.Second * 10)
                return nil
            },
        )

        if err != nil {
            t.Error(err)
        }
    }()

    err := Acquire(
        key, 30, 3, func() error {
            fmt.Println("do something!")
            return nil
        },
    )
    if err != nil {
        t.Fatal(err)
    }
}
