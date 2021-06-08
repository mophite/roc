package output

import (
	"bytes"

	"roc/rlog/common"
	"roc/rlog/output/console"
)

type Outputor interface {
	Init(string)
	Out(level common.Level, b *bytes.Buffer)
	Level() common.Level
	SetLevel(level common.Level)
	Poller()
	Close()
	String() string
}

var DefaultOutput Outputor = &console.Console{}
