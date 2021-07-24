package mock

import (
	"testing"

	"github.com/go-roc/roc/internal/etcd"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := etcd.DefaultEtcd.Put("/test/roc", "data")
	if err != nil {
		t.Fatal(err)
	}
	v, err := etcd.DefaultEtcd.GetWithKey("/test/roc")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(v), "data")

	Stop()
}
