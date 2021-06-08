package console

import (
	"bytes"
	"fmt"

	"roc/internal/x/bytesbuffpool"
	"roc/rlog/common"
)

type Console struct {
	level common.Level
}

func (s *Console) Init(string) {
	return
}

func (s *Console) Out(level common.Level, b *bytes.Buffer) {
	if level < s.level {
		return
	}

	fmt.Printf(b.String())

	bytesbuffpool.Put(b)
}

func (s *Console) Level() common.Level {
	return s.level
}

func (s *Console) SetLevel(l common.Level) {
	s.level = l
}

func (s *Console) Poller() {
	return
}

func (s *Console) Close() {
	return
}

func (s *Console) String() string {
	return "console"
}
