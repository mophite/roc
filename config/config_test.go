package config

import (
    "encoding/json"
    "testing"
    "time"

    _ "github.com/go-roc/roc/etcd/mock"
)

func init() {
    err := NewConfig()
    if err != nil {
        panic(err)
    }
}

type result struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func TestGlobal(t *testing.T) {
    const data = `{ "name":"roc", "age":17 }`
    const testKey = "test"

    err := gRConfig.opts.e.Put(gRConfig.opts.global+testKey, data)
    if err != nil {
        t.Fatal(err)
    }

    //get global config
    var r result
    err = json.Unmarshal(getDataBytes(gRConfig.opts.prefix+".test"), &r)
    if err != nil {
        t.Fatal(err)
    }

    if r.Name != "roc" {
        t.Fatal("not equal")
    }

    d := `{ "name":"roc", "age":18 }`
    //cover global
    err = gRConfig.opts.e.Put(gRConfig.opts.private+testKey, d)
    if err != nil {
        t.Fatal(err)
    }

    err = json.Unmarshal(getDataBytes(gRConfig.opts.prefix+".test"), &r)
    if err != nil {
        t.Fatal(err)
    }

    if r.Age != 18 {
        t.Fatal("not equal")
    }
}

func TestDecode(t *testing.T) {

    const data = `{ "name":"roc", "age":17 }`
    const testKey = "test"

    err := gRConfig.opts.e.Put(gRConfig.opts.global+testKey, data)
    if err != nil {
        t.Fatal(err)
    }

    var r result
    err = Decode2Config(gRConfig.opts.prefix+".test", &r)
    if err != nil {
        t.Fatal(err)
    }
    t.Log(r)

    d := `{ "name":"roc", "age":20 }`
    //cover global
    err = gRConfig.opts.e.Put(gRConfig.opts.global+testKey, d)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(time.Second * 2)
    t.Log(r)

    if r.Age != 20 {
        t.Fatal("not equal")
    }
}
