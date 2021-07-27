package main

import (
    "fmt"

    "github.com/go-roc/roc/config"

    _ "github.com/go-roc/roc/internal/etcd/mock"
)

func main() {

    //new config use default option
    err := config.NewConfig(config.LocalFile())
    if err != nil {
        panic(err)
    }

    const key = "test"
    var result struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    err = config.DecodePrivate(key, &result)
    if err != nil {
        panic(err)
    }

    fmt.Println("--------", result)
}
