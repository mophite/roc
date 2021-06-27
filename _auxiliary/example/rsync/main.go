package main

import (
    "fmt"
    "time"

    _ "github.com/go-roc/roc/internal/etcd/mock"
    "github.com/go-roc/roc/x/rsync"
)

func main() {
	const key = "test"
	go func() {
		err := rsync.Acquire(
			key, 30, 3, func() error {
				fmt.Println("do something!")
				time.Sleep(time.Second * 10)
				return nil
			},
		)

		if err != nil {
			panic(err)
		}
	}()

	err := rsync.Acquire(
		key, 30, 3, func() error {
			fmt.Println("do something!")
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
}
