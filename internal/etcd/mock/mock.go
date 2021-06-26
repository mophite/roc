package mock

import (
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/go-roc/roc/internal/etcd"
)

func init() {
	err := etcd.NewEtcd(
		time.Second*5, 30, &clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	)

	if err != nil {
		panic(err)
	}
}
