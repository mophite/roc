package upload

import (
    "os"

    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
)

type File struct{}

func (h *File) Upload(c *context.Context, req *phello.FileReq, rsp *phello.FileRsp) (err error) {
    f, err := os.OpenFile(req.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
    if err != nil {
        c.Error(err)
        return
    }

    f.Write(req.Body)

    rsp.Msg = "success"
    rsp.Code = 200

    return nil
}
