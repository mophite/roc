package roc

import (
	"os"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/go-roc/roc/config"

	"github.com/go-roc/roc/etcd"
	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/internal/sig"
	"github.com/go-roc/roc/internal/transport"
	rs "github.com/go-roc/roc/internal/transport/rscoket"
	"github.com/go-roc/roc/internal/x"
	"github.com/go-roc/roc/parcel"
	"github.com/go-roc/roc/parcel/codec"
)

const SupportPackageIsVersion1 = 1

type Option struct {
	//server id can be set,or random
	//this is the unique identifier of the service
	id string

	//server name,eg.srv.hello or api.hello
	name string

	//schema is the namespace in your service collection
	//it's the root name prefix
	schema string

	//random port
	//Random ports make you don’t need to care about specific port numbers,
	//which are more commonly used in internal services
	randPort *[2]int

	//socket tcp ip:port address
	tcpAddress string

	//websocket ip:port address
	wssAddress string

	//websocket relative path address of websocket
	wssPath string

	//server version
	version string

	//buffSize to data tunnel if it's need
	buffSize int

	//server transport
	server transport.Server

	//error packet
	//It will affect the format of the data you return
	err parcel.ErrorPackager

	//service discovery endpoint
	e *endpoint.Endpoint

	//receive system signal
	signal []os.Signal

	//wrapper some middleware
	//it's can be interrupt
	wrappers []parcel.Wrapper

	//when server exit,will do exit func
	exit []func()

	// connect server within connectTimeout
	// if out of ranges,will be timeout
	connectTimeout time.Duration

	// keepalive setting,the period for requesting heartbeat to stay connected
	keepaliveInterval time.Duration

	// keepalive setting,the longest time the connection can survive
	keepaliveLifetime time.Duration

	// transport client
	client transport.Client

	//service discover registry
	registry registry.Registry

	//for requestResponse try to retry request
	retry int

	//data encoding or decoding
	cc codec.Codec

	etcdConfig *clientv3.Config

	configOpt []config.Options
}

type Options func(option *Option)

func ConfigOption( opts ...config.Options) Options {
	return func(option *Option) {
		option.configOpt = opts
	}
}

// EtcdConfig setting global etcd config first
func EtcdConfig(e *clientv3.Config) Options {
	return func(option *Option) {
		option.etcdConfig = e
	}
}

func Codec(cc codec.Codec) Options {
	return func(option *Option) {
		option.cc = cc
	}
}

func BuffSize(buffSize int) Options {
	return func(option *Option) {
		option.buffSize = buffSize
	}
}

func Wrapper(wrappers ...parcel.Wrapper) Options {
	return func(option *Option) {
		option.wrappers = append(option.wrappers, wrappers...)
	}
}

func Exit(exit ...func()) Options {
	return func(option *Option) {
		option.exit = exit
	}
}

func Signal(signal ...os.Signal) Options {
	return func(option *Option) {
		option.signal = signal
	}
}

func E(e *endpoint.Endpoint) Options {
	return func(option *Option) {
		option.e = e
	}
}

func Port(port [2]int) Options {
	return func(option *Option) {
		if port[0] > port[1] {
			panic("port index 0 must more than 1")
		}

		if port[0] < 10000 {
			panic("rand port for internal server suggest more than 10000")
		}

		option.randPort = &port
	}
}

func Error(err parcel.ErrorPackager) Options {
	return func(option *Option) {
		option.err = err
	}
}

func WssAddress(address, path string) Options {
	return func(option *Option) {
		option.wssAddress = address
		option.wssPath = path
	}
}

func Id(id string) Options {
	return func(option *Option) {
		option.id = id
	}
}

func Namespace(name string) Options {
	return func(option *Option) {
		option.name = name
	}
}

func TCPAddress(address string) Options {
	return func(option *Option) {
		option.tcpAddress = address
	}
}

func Version(version string) Options {
	return func(option *Option) {
		option.version = version
	}
}

// ConnectTimeout set connect timeout
func ConnectTimeout(connectTimeout time.Duration) Options {
	return func(option *Option) {
		option.connectTimeout = connectTimeout
	}
}

// KeepaliveInterval set keepalive interval
func KeepaliveInterval(keepaliveInterval time.Duration) Options {
	return func(option *Option) {
		option.keepaliveInterval = keepaliveInterval
	}
}

// KeepaliveLifetime set keepalive life time
func KeepaliveLifetime(keepaliveLifetime time.Duration) Options {
	return func(option *Option) {
		option.keepaliveLifetime = keepaliveLifetime
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.etcdConfig == nil {
		opt.etcdConfig = &clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Second * 5,
		}
	}

	config.NewConfig(opt.configOpt...)

	// init etcd.DefaultEtcd
	err := etcd.NewEtcd(time.Second*5, 300, opt.etcdConfig)
	if err != nil {
		panic("etcdConfig occur error:" + err.Error())
	}

	if opt.name == "" {
		opt.name = x.GetProjectName()
	}

	if opt.id == "" {
		//todo change to git commit id+timestamp
		opt.id = x.NewUUID()
	}

	ip, err := x.LocalIp()
	if err != nil {
		panic(err)
	}

	if opt.randPort == nil {
		opt.randPort = &[2]int{10000, 19999}
	}

	if opt.tcpAddress == "" {
		opt.tcpAddress = ip + ":" + strconv.Itoa(x.RandInt(opt.randPort[0], opt.randPort[1]))
	}

	if opt.version == "" {
		opt.version = namespace.DefaultVersion
	}

	if opt.wssAddress != "" && opt.wssPath == "" {
		opt.wssPath = "/roc/wss"
	}

	//if opt.ratelimit <= 0 {
	//	opt.ratelimit = math.MaxInt32
	//}

	if opt.schema == "" {
		opt.schema = namespace.DefaultSchema
	}

	if opt.err == nil {
		opt.err = parcel.DefaultErrorPacket
	}

	if opt.server == nil {
		opt.server = rs.NewServer(
			opt.tcpAddress,
			opt.wssAddress,
			opt.name,
			opt.buffSize,
		)
	}

	opt.registry = registry.NewRegistry(registry.Schema(opt.schema))

	if opt.e == nil {
		opt.e = &endpoint.Endpoint{
			Id:      opt.id,
			Name:    opt.name,
			Version: opt.version,
			Address: opt.tcpAddress,
		}

		opt.e.Splicing(opt.schema)
	}

	if opt.signal == nil {
		opt.signal = sig.DefaultSignal
	}

	if opt.buffSize == 0 {
		opt.buffSize = 10
	}

	if opt.cc == nil {
		opt.cc = codec.DefaultCodec
	}

	//set connect timeout or default
	if opt.connectTimeout <= 0 {
		opt.connectTimeout = time.Second * 5
	}

	if opt.keepaliveLifetime <= 0 {
		opt.keepaliveLifetime = time.Second * 600
	}

	if opt.keepaliveInterval <= 0 {
		opt.keepaliveInterval = time.Second * 5
	}

	if opt.client == nil {
		//default is rsocket
		opt.client = rs.NewClient(
			opt.connectTimeout,
			opt.keepaliveInterval,
			opt.keepaliveLifetime,
		)
	}

	return opt
}
