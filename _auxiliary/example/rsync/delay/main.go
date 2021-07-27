package main

import (
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
		err := rsync.AcquireDelay(
			key, 5, func() error {
				c <- "c1"
				fmt.Println("do something one!")
				return nil
			},
		)
		fmt.Println("-------one err", nil)

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := rsync.AcquireDelay(
			key, 5, func() error {
				time.Sleep(time.Second * 1)
				c <- "c2"
				fmt.Println("do something two!")
				return nil
			},
		)

		fmt.Println("-----two err----",err)

		if err != nil {
			panic(err)
		}
	}()

	<-c
	<-c

	time.Sleep(time.Second)
}
