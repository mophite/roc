module roc

go 1.16

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/go-roc/roc v0.0.8
	github.com/google/uuid v1.2.0
	github.com/jjeffcaii/reactor-go v0.5.1
	github.com/json-iterator/go v1.1.10
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/rsocket/rsocket-go v0.8.4
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.36.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
