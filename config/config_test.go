package config

import (
	"testing"

	"github.com/go-roc/roc/internal/etcd/mock"
	_ "github.com/go-roc/roc/internal/etcd/mock"
	"github.com/stretchr/testify/assert"
)

type configData struct {
	Name string
	Age  int
}

func setup() {
	err := NewConfig(
		Public("public"),
		Private("private"),
		DisableDynamic(),
		LogOut(),
		Prefix("roc."),
	)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	mock.Stop()
}

func TestDecodePublic(t *testing.T) {
	setup()

	defer teardown()

	var err error

	data := `{"name":"roc","age":1}`
	key := "test"

	//key must equal
	assert.Equal(t, "configroc/v1.0.0/public/roc.test", gRConfig.opts.public+gRConfig.opts.publicPrefix+key)

	err = PutPublic(key, data)
	if err != nil {
		t.Fatal(err)
	}

	var d configData
	err = DecodePublic(key, &d)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, d.Name, "roc")
	assert.Equal(t, d.Age, 1)

	var d1 configData
	err = DecodePrivate(key, &d1)
	assert.NotNil(t, err)
	assert.Equal(t, d1.Name, "")
	assert.Equal(t, d1.Age, 0)
}

func TestDecodePrivate(t *testing.T) {
	setup()

	defer teardown()

	var err error

	data := `{"name":"roc","age":1}`
	key := "test"

	//key must equal
	assert.Equal(t, "configroc/v1.0.0/private/test", gRConfig.opts.private+key)

	err = PutPrivate(key, data)
	if err != nil {
		t.Fatal(err)
	}

	var v1 configData
	err = DecodePrivate(key, &v1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, v1.Name, "roc")
	assert.Equal(t, v1.Age, 1)

	var v2 configData
	err = DecodePublic(key, &v2)
	assert.NotNil(t, err)
	assert.Equal(t, v2.Name, "")
	assert.Equal(t, v2.Age, 0)
}
