package ipc

import (
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/service"
)

func InitIpc(s *service.Service) {
    sayClient = phello.NewHelloSrvClient(s.Client())
    imClient = phello.NewImClient(s.Client())
}
