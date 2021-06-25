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
	var testKey = "test"

	//key: /public/roc.test
	err := gRConfig.opts.e.Put(gRConfig.opts.public+gRConfig.opts.prefix+testKey, data)
	if err != nil {
		t.Fatal(err)
	}

	//get public config
	var r result
	//key roc.test
	err = json.Unmarshal(getDataBytes(gRConfig.opts.prefix+"test"), &r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Name != "roc" {
		t.Fatal("not equal")
	}

	var r1 result
	d := `{ "name":"roc", "age":18 }`
	//cover public
	//key: /private/roc.test
	err = gRConfig.opts.e.Put(gRConfig.opts.private+gRConfig.opts.prefix+testKey, d)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	//key: roc.test
	err = json.Unmarshal(getDataBytes(gRConfig.opts.prefix+"test"), &r1)
	if err != nil {
		t.Fatal(err)
	}

	if r1.Age != 18 {
		t.Fatal("not equal")
	}
}

func TestDecode(t *testing.T) {
	const data = `{ "name":"roc", "age":17 }`
	var testKey = "test"

	//key: /public/roc.test
	err := gRConfig.opts.e.Put(gRConfig.opts.public+gRConfig.opts.prefix+testKey, data)
	if err != nil {
		t.Fatal(err)
	}

	//get public config
	var r result
	//key roc.test
	err = Decode2Config(gRConfig.opts.prefix+"test", &r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Name != "roc" {
		t.Fatal("not equal")
	}

	d := `{ "name":"roc", "age":18 }`
	//cover public
	//key: /private/roc.test
	err = gRConfig.opts.e.Put(gRConfig.opts.private+gRConfig.opts.prefix+testKey, d)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	if r.Age != 18 {
		t.Fatal("not equal")
	}
}
