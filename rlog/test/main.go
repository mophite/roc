package main

import (
	"time"

	"roc/rlog"
)

func main() {
	rlog.Infof("------%v", 111)
	time.Sleep(time.Second * 2)
}
