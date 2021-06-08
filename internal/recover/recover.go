package recover

import (
	"roc/rlog"
)

func Recover(f ...func()) {
	for i := range f {
		f[i]()
	}

	if err := recover(); err != nil {
		rlog.Stack(err)
	}
}
