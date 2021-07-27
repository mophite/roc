package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-roc/roc/internal/etcd/mock"
	"github.com/go-roc/roc/x/rsync"
)

func main() {
	const key = "test"
	var c = make(chan string)
	go func() {
		time.Sleep(time.Second * 2)
		err := rsync.AcquireOnce(
			key, 5, func() error {
				c <- "c1"
				fmt.Println("do something one!")
				return nil
			},
		)

		if err != nil && err == context.DeadlineExceeded {
			fmt.Println("-------one err", err)
			c <- "c1"
		}
	}()

	go func() {
		err := rsync.AcquireOnce(
			key, 5, func() error {
				time.Sleep(time.Second * 1)
				c <- "c2"
				fmt.Println("do something two!")
				return nil
			},
		)

		if err != nil &&err==context.DeadlineExceeded{
			fmt.Println("-----two err----", err)
		}
	}()

	<-c
	<-c

	time.Sleep(time.Second)
}
