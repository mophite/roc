package ipc

import (
    "github.com/go-roc/roc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
)

func InitIpc(s *roc.Service) {
    sayClient = phello.NewHelloSrvClient(s.Client())
    imClient = phello.NewImClient(s.Client())
}
