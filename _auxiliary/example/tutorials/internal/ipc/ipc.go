package ipc

import (
    "github.com/go-roc/roc"
)

type Ipc interface {
    Service() *roc.Service
    InvokeOptions() []roc.InvokeOptions
}
