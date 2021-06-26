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

	simple(key, &result)
	coverPublic(key, &result)
}

//put key/value to etcd:
//go:generate etcdctl configroc/v1.0.0/public/roc.test { "name":"roc", "age":17 }
func simple(key string, v interface{}) {
	//simple public use
	//the key is roc.test
	err := config.DecodePublic(key, v)
	if err != nil {
		panic(err)
	}

	fmt.Println("------", v)
	//output: ------ {roc 17}
}

//put key/value to etcd:
//go:generate etcdctl configroc/v1.0.0/private/roc.test { "name":"roc", "age":18 }
func coverPublic(key string, v interface{}) {
	//the key is roc.test
	//cover public by private
	err := config.DecodePublic(key, v)
	if err != nil {
		panic(err)
	}

	fmt.Println("------", v)
	//output: ------ {roc 18}
}
