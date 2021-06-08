package sig

import (
	"os"
	"syscall"
)

var DefaultSignal = []os.Signal{
	syscall.SIGINT,
	syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGKILL,
	syscall.SIGQUIT,
}

func ExitSignal(signal ...os.Signal) []os.Signal {
	return signal
}
