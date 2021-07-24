package roc

import (
    "net/http"
    "strings"

    "github.com/go-roc/roc/service"
    "github.com/gorilla/mux"
)

const SupportPackageIsVersion1 = 1
const HttpHandlePrefix = "RocApi"

func NewRocService() *service.Service{

}

var defaultRouter *mux.Router

// GROUP for not post http method
func GROUP(prefix string) {
    if !strings.HasPrefix(prefix, "/") {
        prefix = "/" + prefix
    }
    if !strings.HasSuffix(prefix, "/") {
        prefix = prefix + "/"
    }

    //cannot had post prefix
    if strings.HasPrefix(prefix, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }

    defaultRouter = service.DefaultService.PathPrefix(prefix).Subrouter()
}

func GET(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodOptions, http.MethodGet)
}

func POST(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPost)
}

func PUT(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPut)
}

func DELETE(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodDelete)
}

func ANY(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler)
}

func HEAD(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodHead)
}

func PATCH(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPatch)
}

func CONNECT(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodConnect)
}

func TRACE(relativePath string, handler http.Handler) {

    relativePath = tidyRelativePath(relativePath)

    if strings.HasPrefix(relativePath, service.DefaultService.GetApiPrefix()) {
        panic("cannot contain unique prefix")
    }
    defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodTrace)
}

func tidyRelativePath(relativePath string) string {
    //trim suffix "/"
    if strings.HasSuffix(relativePath, "/") {
        relativePath = strings.TrimSuffix(relativePath, "/")
    }

    //add prefix "/"
    if !strings.HasPrefix(relativePath, "/") {
        relativePath = "/" + relativePath
    }

    return relativePath
}
