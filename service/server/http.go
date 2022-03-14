package server

import (
    "bytes"
    "io"
    "net/http"
    "strings"

    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/parcel"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/parcel/packet"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service/router"
)

func handlerServerHttp(c *context.Context, s *Server, w http.ResponseWriter, r *http.Request) {
    c.Metadata.SetMethod(r.URL.Path)

    c.Writer = w
    c.Request = r

    for k, v := range r.Header {
        if len(v) == 0 {
            continue
        }
        c.SetHeader(k, v[0])
    }
    c.ContentType = c.GetHeader(namespace.DefaultHeaderContentType)
    c.SetCodec()

    c.RemoteAddr = r.RemoteAddr

    w.Header().Set(namespace.DefaultHeaderContentType, c.ContentType)
    w.Header().Set(namespace.DefaultHeaderTrace, c.Trace.TraceId())

    for i := range s.opts.Dog {
        rsp, err := s.opts.Dog[i](c)
        if err != nil {
            c.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(c.Codec().MustEncode(rsp))
            return
        }
    }

    switch r.Method {
    case http.MethodPost, http.MethodDelete:
        var req, rsp = parcel.PayloadIo(r.Body), parcel.NewPacket()

        _ = r.Body.Close()

        methodPostOrDelete(c, s, w, r, req, rsp)
        parcel.Recycle(req)
        parcel.Recycle(rsp)

        return

    case http.MethodPut:
        f, h, err := r.FormFile("file")
        if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(s.opts.Err.Error400(c))
            return
        }

        var buf = bytes.NewBuffer(make([]byte, 0, 10485760))

        io.Copy(buf, f)

        var fileReq = &packet.FileReq{}
        fileReq.Body = buf.Bytes()
        fileReq.FileSize = h.Size
        fileReq.FileName = h.Filename
        fileReq.Extra = r.FormValue("extra")

        fb, err := c.Codec().Encode(fileReq)
        if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(s.opts.Err.Error400(c))
            return
        }

        var req, rsp = parcel.Payload(fb), parcel.NewPacket()

        methodPut(c, s, w, r, req, rsp)
        _ = r.Body.Close()

        parcel.Recycle(req)
        parcel.Recycle(rsp)

        return

    case http.MethodGet:

        if strings.Count(r.URL.Path, "/") == 1 {
            c.Metadata.M = s.opts.RootRouterRedirect
        }

        values := r.URL.Query()

        var apiReq = &packet.ApiReq{Params: make(map[string]string, len(values))}
        for k, v := range values {
            if len(v) > 0 {
                apiReq.Params[k] = v[0]
            }
        }

        fb, err := c.Codec().Encode(apiReq)
        if err != nil {
            rlog.Error(err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write(s.opts.Err.Error400(c))
            return
        }

        var req, rsp = parcel.Payload(fb), parcel.NewPacket()

        methodGet(c, s, w, r, req, rsp)

        parcel.Recycle(req)
        parcel.Recycle(rsp)

        return

    case http.MethodOptions:
        w.WriteHeader(http.StatusOK)
        return
    }

    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write(s.opts.Err.Error405(c))
}

func methodPostOrDelete(
    c *context.Context,
    s *Server,
    w http.ResponseWriter,
    r *http.Request,
    req, rsp *parcel.RocPacket,
) {

    err := s.route.RRProcess(c, req, rsp)

    if err == router.ErrNotFoundHandler {
        c.Errorf("path=%s |service=%s", r.URL.Path, s.opts.Name)
        w.WriteHeader(http.StatusNotFound)
        w.Write(s.opts.Err.Error404(c))
        return
    }

    if len(rsp.Bytes()) > 0 {
        w.WriteHeader(http.StatusOK)
        w.Write(rsp.Bytes())
    } else if err != nil {
        rlog.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(s.opts.Err.Error500(c))
    }
}

func methodPut(
    c *context.Context,
    s *Server,
    w http.ResponseWriter,
    r *http.Request,
    req, rsp *parcel.RocPacket,
) {

    c.IsPutFile = true
    err := s.route.RRProcess(c, req, rsp)

    if err == router.ErrNotFoundHandler {
        c.Errorf("err=%v |path=%s", err, c.Metadata.Method())
        w.Header().Set("Content-type", "text/plain")
        w.WriteHeader(http.StatusNotFound)
        w.Write(s.opts.Err.Error404(c))
        return
    }

    if len(rsp.Bytes()) > 0 {
        w.WriteHeader(http.StatusOK)
        w.Write(rsp.Bytes())
    } else if err != nil {
        rlog.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(s.opts.Err.Error500(c))
    }
}

func methodGet(
    c *context.Context,
    s *Server,
    w http.ResponseWriter,
    r *http.Request,
    req, rsp *parcel.RocPacket,
) {
    err := s.route.RRProcess(c, req, rsp)

    if err == router.ErrNotFoundHandler {
        c.Errorf("err=%v |path=%s", err, c.Metadata.Method())
        w.WriteHeader(http.StatusNotFound)
        w.Write(s.opts.Err.Error404(c))
        return
    }

    if len(rsp.Bytes()) > 0 {
        w.WriteHeader(http.StatusOK)
        w.Write(rsp.Bytes())
    } else if err != nil {
        rlog.Error(err)
        w.Header().Set("Content-type", "text/plain")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(s.opts.Err.Error500(c))
    }
}
