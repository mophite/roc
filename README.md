# Roc

![logo](./logo.jpg)

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

```go
    go env -w GO111MODULE=on
```
```go
    go get github.com/go-roc/roc/cmd/protoc-gen-roc
```

- generate proto file to go
  file,like [hello.proto](https://roc/_auxiliary/example/tutorials/proto/pbhello.proto)

```go
    protoc --roc_out = plugins = roc:.*.proto
```

- run a roc server

```go
    var s = server.NewRocServer(server.Namespace("srv.hello"))
    pbhello.RegisterHelloWorldServer(s, &Hello{})
    err := s.Run()
```

- client rpc to server

```go
    var opt = client.WithScope("srv.hello")
    var client = pbhello.NewHelloWorldClient(client.NewRocClient())
    rsp, err := h.client.Say(context.Background(), &pbhello.SayReq{Inc: 1}, h.opt)
```

### 💞️ see more [example](https://roc/tree/master/_auxiliary/example) for more help.

### 📫 How to reach me by email ...
```email
  1743299@qq.com
```

![code](./qr.png)

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



