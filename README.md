# Roc

![logo](https://github.com/go-roc/roc/blob/master/_auxiliary/imgs/logo.jpg)

![GitHub Workflow Status](https://github.com/rsocket/rsocket-go/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-roc/roc)](https://goreportcard.com/report/github.com/go-roc/roc)
[![Go Reference](https://pkg.go.dev/badge/roc.svg)](https://pkg.go.dev/github.com/go-roc/roc)
![GitHub](https://img.shields.io/github/license/go-roc/roc?logo=rsocket)
![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/go-roc/roc?include_prereleases)

### 👋 Roc is a rpc micro framework,it designed with go,and transport protocol by [rsocket-go](https://github.com/rsocket/rsocket-go).

<br>***IT IS UNDER ACTIVE DEVELOPMENT, APIs are unstable and maybe change at any time until release of v1.0.0.***

### 👀 Features

- Simple to use ✨
- Lightweight ✨
- High performance ✨
- Service discovery ✨
- Support websocket and socket same time ✨
- Support json or [gogo proto](https://github.com/gogo/protobuf) when use rpc ✨

### 🌱 Quick start

- first you must install [proto](https://github.com/gogo/protobuf) and [etcd](https://github.com/etcd-io/etcd)

- install protoc-gen-roc

```shell
    GO111MODULE=on go get github.com/go-roc/roc/cmd/protoc-gen-roc
```

- generate proto file to go file,like [hello.proto](https://roc/_auxiliary/example/tutorials/proto/pbhello.proto)

```shell
    protoc --roc_out = plugins = roc:.*.proto
```

- run a roc service

```go
package main

import (
    "github.com/go-roc/roc/_auxiliary/example/hello/app/api/api.hello/hello"
    "github.com/go-roc/roc/_auxiliary/example/hello/proto/phello"
    "github.com/go-roc/roc/rlog"
    "github.com/go-roc/roc/service"
)

func main() {
    s := service.New(
        service.HttpAddress("0.0.0.0:9999"),
        service.Namespace("api.hello"),
        service.TCPAddress("0.0.0.0:8888"),
        service.WssAddress("0.0.0.0:10000", "/hello"),
    )

    phello.RegisterHelloWorldServer(s.Server(), &hello.Hello{Client: s.Client()})
    err := s.Run()
    if err != nil {
        rlog.Error(err)
    }
}
```

- config help

```go
package main

import (
    "fmt"

    "github.com/go-roc/roc/config"

    _ "github.com/go-roc/roc/internal/etcd/mock"
)

func main() {

    //new config use default option
    err := config.NewConfig()
    if err != nil {
        panic(err)
    }

    const key = "test"
    var result struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    coverPublic(key, &result)
    coverPrivate(key, &result)
}

//put key/value to etcd:
//go:generate etcdctl put  configroc/v1.0.0/public/roc.test "{ "name":"roc", "age":18 }"
func coverPublic(key string, v interface{}) {
    //simple public use
    //the key is roc.test
    err := config.DecodePublic(key, v)
    if err != nil {
        panic(err)
    }

    fmt.Println("------", v)
    //output: ------ {roc 18}
}

//put key/value to etcd:
//go:generate etcdctl put  configroc/v1.0.0/private/test "{ "name":"roc", "age":18 }"
func coverPrivate(key string, v interface{}) {
    //the key is test
    err := config.DecodePrivate(key, v)
    if err != nil {
        panic(err)
    }

    fmt.Println("------", v)
    //output: ------ {roc 18}
}

```

### 💞️ see more [example](https://github.com/go-roc/roc/tree/master/_auxiliary/example) for more help.

### 📫 How to reach me and be a contributor ...

### ✨ TODO ✨

- [ ] bench test
- [ ] sidecar
- [ ] more example
- [ ] more singleton tests
- [ ] generate dir
- [ ] command for request service
- [ ] sidecar service
- [ ] config service
- [ ] broker redirect request service
- [ ] logger service
- [ ] simple service GUI



