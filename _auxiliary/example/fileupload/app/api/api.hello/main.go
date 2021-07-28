package main

import (
	"github.com/go-roc/roc/_auxiliary/example/fileupload/app/api/api.hello/upload"
	"github.com/go-roc/roc/_auxiliary/example/fileupload/proto/phello"
	"github.com/go-roc/roc/rlog"
	"github.com/go-roc/roc/service"
)

func main() {
	s := service.New(
		service.Namespace("api.hello"),
		service.HttpAddress("0.0.0.0:9999"),
		service.TCPAddress("0.0.0.0:8888"),
	)

	phello.RegisterFileServer(s.Server(), &upload.File{})

	err := s.Run()
	if err != nil {
		rlog.Error(err)
	}
}
