module protoc-gen-roc

go 1.16

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/go-roc/roc v0.9.2
	github.com/gogo/protobuf v1.3.2
)
